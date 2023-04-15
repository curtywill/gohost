package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gohost/requests"
	"gohost/structs"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

type Markdown struct {
	block structs.MarkdownBlock
}

func NewMarkdown(content string) Markdown {
	block := structs.MarkdownBlock{Content: content}
	return Markdown{block}
}

func (m Markdown) GetBlock() structs.Blocks {
	return structs.Blocks{
		Type:     "markdown",
		Markdown: &m.block,
	}
}

type Attachment struct {
	block         structs.AttachmentBlock
	filepath      string
	filename      string
	contentType   string
	contentLength string
}

func NewAttachment(filepath, altText string) Attachment {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}
	filename := stat.Name()
	sl := strings.Split(filename, ".")
	contentType := mime.TypeByExtension("." + strings.ToLower(sl[len(sl)-1]))
	contentLength := fmt.Sprint(stat.Size())

	block := structs.AttachmentBlock{
		AttachmentId: "00000000-0000-0000-0000-000000000000",
		AltText:      altText,
	}

	return Attachment{block, filepath, filename, contentType, contentLength}
}

func (a Attachment) GetBlock() structs.Blocks {
	return structs.Blocks{
		Type:       "attachment",
		Attachment: &a.block,
	}
}

func (a *Attachment) upload(postId int, project Project) {
	form := makeForm(a.filename, a.contentType, a.contentLength)

	formHeader := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}
	r := structs.AttachStart{}
	data, _ := requests.Fetch("POST",
		fmt.Sprintf("/project/%s/posts/%d/attach/start", project.Handle(), postId),
		project.u.cookie, formHeader, strings.NewReader(form), false)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(a.filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	contentTypeHeader, body, err := doSpacesForm(r.RequiredFields, a.filename, file)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, r.Url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", contentTypeHeader)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 204 {
		log.Fatalf("bad status code from redcent dev: %d", res.StatusCode)
	}
	requests.Fetch("POST",
		fmt.Sprintf("/project/%s/posts/%d/attach/finish/%s", project.Handle(), postId, r.AttachmentId),
		project.u.cookie, nil, nil, false)

	// file's been attached, need to add attachment ID to attachment struct
	a.block.AttachmentId = r.AttachmentId
}

// helper functions for Upload
func makeForm(filename, contentType, contentLength string) string {
	formData := url.Values{}
	formData.Set("content_length", contentLength)
	formData.Set("content_type", contentType)
	formData.Set("filename", filename)
	encodedForm := formData.Encode()
	return encodedForm
}

func doSpacesForm(r structs.RequiredFields, filename string, file *os.File) (string, *bytes.Buffer, error) {
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
		fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename)) // TODO: escape filename
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
