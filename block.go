package gohost

// Cohost posts are comprised of blocks, this file contains the
// structs that symbolize said blocks
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"github.com/curtywill/gohost/requests"
)

type Markdown struct {
	block markdownBlockResponse
}

// Returns a markdown block that represent text data on a Cohost post
func MarkdownBlock(content string) Markdown {
	block := markdownBlockResponse{Content: content}
	return Markdown{block}
}

func (m Markdown) getBlock() blocksResponse {
	return blocksResponse{
		Type:     "markdown",
		Markdown: &m.block,
	}
}

type Attachment struct {
	block         attachmentBlockResponse
	filepath      string
	filename      string
	contentType   string
	contentLength string
}

// Returns an Attachment struct that contains information about an image given a filepath.
// Can return a path error.
func AttachmentBlock(filepath, altText string) (Attachment, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Attachment{}, err
	}
	defer file.Close()

	stat, err := os.Stat(filepath)
	if err != nil {
		return Attachment{}, err
	}
	filename := stat.Name()
	sl := strings.Split(filename, ".")
	contentType := mime.TypeByExtension("." + strings.ToLower(sl[len(sl)-1]))
	contentLength := fmt.Sprint(stat.Size())

	block := attachmentBlockResponse{
		AltText: altText,
	}

	return Attachment{block, filepath, filename, contentType, contentLength}, nil
}

func (a Attachment) getBlock() blocksResponse {
	if a.block.AttachmentId == "" {
		a.block.AttachmentId = "00000000-0000-0000-0000-000000000000"
	}
	return blocksResponse{
		Type:       "attachment",
		Attachment: &a.block,
	}
}

func (a *Attachment) upload(client *http.Client, postId int, project Project) error {
	if a.block.AttachmentId != "" {
		return nil
	}

	form := makeForm(a.filename, a.contentType, a.contentLength)

	formHeader := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	r := attachStartResponse{}
	data, _, err := requests.Fetch(client, "POST",
		fmt.Sprintf("/project/%s/posts/%d/attach/start", project.Handle(), postId),
		project.u.cookie, formHeader, strings.NewReader(form), false)

	if err != nil {
		return err
	}

	json.Unmarshal(data, &r)

	file, err := os.Open(a.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	contentTypeHeader, body, err := doSpacesForm(r.RequiredFields, a.filename, file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, r.Url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentTypeHeader)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("error uploading attachment to digital ocean! status code: %d", res.StatusCode)
	}

	_, _, err = requests.Fetch(client, "POST",
		fmt.Sprintf("/project/%s/posts/%d/attach/finish/%s", project.Handle(), postId, r.AttachmentId),
		project.u.cookie, nil, nil, false)

	if err != nil {
		return err
	}
	// file's been attached, need to add attachment ID to attachment struct
	a.block.AttachmentId = r.AttachmentId

	return nil
}

// helper functions for Upload

// creates the url encoded form that starts the attachment process
func makeForm(filename, contentType, contentLength string) string {
	formData := url.Values{}
	formData.Set("content_length", contentLength)
	formData.Set("content_type", contentType)
	formData.Set("filename", filename)
	encodedForm := formData.Encode()
	return encodedForm
}

// creates the multipart form needed for the digital ocean spaces credentials
func doSpacesForm(r requiredFieldsResponse, filename string, file *os.File) (string, *bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close() // writes closing boundary line

	writer.WriteField("acl", r.Acl)
	writer.WriteField("Content-Type", r.ContentType)
	writer.WriteField("Content-Disposition", r.ContentDisposition)
	writer.WriteField("Cache-Control", r.CacheControl)
	writer.WriteField("key", r.Key)
	writer.WriteField("bucket", r.Bucket)
	writer.WriteField("X-Amz-Algorithm", r.XAmzAlgorithm)
	writer.WriteField("X-Amz-Credential", r.XAmzCredential)
	writer.WriteField("X-Amz-Date", r.XAmzDate)
	writer.WriteField("Policy", r.Policy)
	writer.WriteField("X-Amz-Signature", r.XAmzSignature)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(filename)))
	h.Set("Content-Type", r.ContentType)
	p, err := writer.CreatePart(h)
	if err != nil {
		return "", nil, err
	}
	_, err = io.Copy(p, file)
	if err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

// helpers for the multipart form
// taken from go source
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
