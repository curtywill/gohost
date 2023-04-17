package models

// structs used for JSON decoding!
// named after the endpoint that the json stream is coming from

type loggedInResponse struct {
	Activated bool   `json:"activated"`
	ModMode   bool   `json:"modMode"`
	ReadOnly  bool   `json:"readOnly"`
	Email     string `json:"email"`
	ProjectId int    `json:"projectId"`
	UserId    int    `json:"userId"`
}

type listEditedProjectsResponse struct {
	Projects []editedProjectResponse `json:"projects"`
}

// TODO: wtf is a flag
type editedProjectResponse struct {
	AvatarShape        string   `json:"avatarShape"`
	AvatarURL          string   `json:"avatarURL"`
	HeaderURL          string   `json:"headerURL"`
	Pronouns           string   `json:"pronouns"`
	URL                string   `json:"url"`
	Privacy            string   `json:"privacy"`
	Dek                string   `json:"dek"`
	Description        string   `json:"description"`
	DisplayName        string   `json:"displayName"`
	Handle             string   `json:"handle"`
	ProjectId          int      `json:"projectId"`
	Flags              [][]byte `json:"flags"`
	FrequentlyUsedTags []string `json:"frequentlyUsedTags"`
}

type projectPostsResponse struct {
	NumItems int             `json:"nItems"`
	Items    []itemsResponse `json:"items"`
}

// add rest of stuff
type itemsResponse struct {
	PostId                   int                   `json:"postId"`
	Headline                 string                `json:"headline"`
	PublishedAt              string                `json:"publishedAt"`
	Filename                 string                `json:"filename"`
	State                    int                   `json:"state"`
	NumComments              int                   `json:"numComments"`
	NumSharedComments        int                   `json:"numSharedComments"`
	ContentWarnings          []string              `json:"cws"`
	Tags                     []string              `json:"tags"`
	Blocks                   []blocksResponse      `json:"blocks"`
	PlainTextBody            string                `json:"plainTextBody"`
	PostingProject           editedProjectResponse `json:"postingProject"`
	TransparentShareOfPostId int                   `json:"transparentShareOfPostId"`
	ShareTree                []itemsResponse       `json:"shareTree"`
	Url                      string                `json:"singlePostPageUrl"`
}

// blocks are pointers to get an empty value (nil in this case) for omitempty to work
type blocksResponse struct {
	Type       string                   `json:"type"`
	Attachment *attachmentBlockResponse `json:"attachment,omitempty"`
	Markdown   *markdownBlockResponse   `json:"markdown,omitempty"`
}

type attachmentBlockResponse struct {
	FileURL      string `json:"fileURL,omitempty"`
	PreviewURL   string `json:"previewURL,omitempty"`
	AttachmentId string `json:"attachmentId"`
	AltText      string `json:"altText"`
}

type markdownBlockResponse struct {
	Content string `json:"content"`
}

type attachStartResponse struct {
	AttachmentId   string                 `json:"attachmentId"`
	Url            string                 `json:"url"`
	RequiredFields requiredFieldsResponse `json:"requiredFields"`
}

type requiredFieldsResponse struct {
	Acl                string `json:"acl"`
	ContentType        string `json:"Content-Type"`
	ContentDisposition string `json:"Content-Disposition"`
	CacheControl       string `json:"Cache-Control"`
	Key                string `json:"key"`
	Bucket             string `json:"bucket"`
	XAmzAlgorithm      string `json:"X-Amz-Algorithm"`
	XAmzCredential     string `json:"X-Amz-Credential"`
	XAmzDate           string `json:"X-Amz-Date"`
	Policy             string `json:"Policy"`
	XAmzSignature      string `json:"X-Amz-Signature"`
}

type postIdResponse struct {
	PostId int `json:"postId"`
}

type saltResponse struct {
	Salt string `json:"salt"`
}
