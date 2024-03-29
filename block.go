package gohost

// Cohost posts are comprised of blocks, this file contains the
// structs that symbolize said blocks
import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"github.com/curtywill/gohost/requests"
)

type Markdown struct {
	Content string
}

func (m Markdown) getBlock() blocksResponse {
	// extracting content like this to suppress "should convert instead of using struct literal"
	// i know  its stupid go leave me alone bxtch
	content := m.Content
	return blocksResponse{
		Type:     "markdown",
		Markdown: &markdownBlockResponse{content},
	}
}

type Attachment struct {
	Filepath      string
	AltText       string
	block         attachmentBlockResponse
	filename      string
	contentType   string
	contentLength int64
	height        int
	width         int
}

func (a Attachment) FileUrl() string {
	return a.block.FileURL
}

func (a *Attachment) initAttachment() error {
	file, err := os.Open(a.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := os.Stat(a.Filepath)
	if err != nil {
		return err
	}

	filename := stat.Name()
	a.filename = filename
	sl := strings.Split(filename, ".")
	a.contentType = mime.TypeByExtension("." + strings.ToLower(sl[len(sl)-1]))
	a.contentLength = stat.Size()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}
	a.height = img.Height
	a.width = img.Width

	a.block = attachmentBlockResponse{
		AltText: a.AltText,
	}

	return nil
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
	// fill Attachment struct before we start uploading
	err := a.initAttachment()
	if err != nil {
		return err
	}

	form := makeForm(a.filename, a.contentType, a.contentLength, a.height, a.width)

	formHeader := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
	}

	r := attachStartResponse{}
	data, _, err := requests.Fetch(client, "POST",
		fmt.Sprintf("/project/%s/posts/%d/attach/start", project.Handle(), postId),
		project.u.cookie, formHeader, bytes.NewReader(form), false)

	if err != nil {
		return err
	}
	json.Unmarshal(data, &r)

	file, err := os.Open(a.Filepath)
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

// creates the form that starts the attachment process
func makeForm(filename, contentType string, contentLength int64, height, width int) []byte {
	// literally a subset of Attachment because I didn't want to export the fields in the Attachment struct
	type form struct {
		Filename      string `json:"filename"`
		ContentType   string `json:"content_type"`
		ContentLength int64  `json:"content_length"`
		Height        int    `json:"height"`
		Width         int    `json:"width"`
	}
	f := form{filename, contentType, contentLength, height, width}
	marshalledForm, _ := json.Marshal(f)

	return marshalledForm
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
