package structs

// structs used for JSON decoding!
// named after the endpoint that the json stream is coming from

type JsonStruct interface {
	LoggedIn |
		ListEditedProjects |
		ProjectPosts |
		AttachStart |
		PostIdStruct |
		Filler
}

type LoggedIn struct {
	Activated bool   `json:"activated"`
	ModMode   bool   `json:"modMode"`
	ReadOnly  bool   `json:"readOnly"`
	Email     string `json:"email"`
	ProjectId int    `json:"projectId"`
	UserId    int    `json:"userId"`
}

type ListEditedProjects struct {
	Projects []EditedProject `json:"projects"`
}

// TODO: investigate real types of flags and frequently use tags (probably strings)
type EditedProject struct {
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
	FrequentlyUsedTags [][]byte `json:"frequentlyUsedTags"`
}

type ProjectPosts struct {
	NumItems int     `json:"nItems"`
	Items    []Items `json:"items"`
}

// add rest of stuff
type Items struct {
	PostId                   int           `json:"postId"`
	Headline                 string        `json:"headline"`
	PublishedAt              string        `json:"publishedAt"`
	Filename                 string        `json:"filename"`
	State                    int           `json:"state"`
	NumComments              int           `json:"numComments"`
	NumSharedComments        int           `json:"numSharedComments"`
	ContentWarnings          []string      `json:"cws"`
	Tags                     []string      `json:"tags"`
	Blocks                   []Blocks      `json:"blocks"`
	PlainTextBody            string        `json:"plainTextBody"`
	PostingProject           EditedProject `json:"postingProject"`
	TransparentShareOfPostId int           `json:"transparentShareOfPostId"`
	ShareTree                []Items       `json:"shareTree"`
}

// blocks are pointers to get an empty value (nil in this case) for omitempty to work
type Blocks struct {
	Type       string           `json:"type"`
	Attachment *AttachmentBlock `json:"attachment,omitempty"`
	Markdown   *MarkdownBlock   `json:"markdown,omitempty"`
}

type AttachmentBlock struct {
	FileURL      string `json:"fileURL,omitempty"`
	PreviewURL   string `json:"previewURL,omitempty"`
	AttachmentId string `json:"attachmentId"`
	AltText      string `json:"altText"`
}

type MarkdownBlock struct {
	Content string `json:"content"`
}

type AttachStart struct {
	AttachmentId   string         `json:"attachmentId"`
	Url            string         `json:"url"`
	RequiredFields RequiredFields `json:"requiredFields"`
}

type RequiredFields struct {
	Acl                string `json:"acl"`
	ContentType        string `json:"Content-Type"`
	ContentDisposition string `json:"Content-Disposition"`
	CacheControl       string `json:"Cache-Control"`
	Key                string `json:"key"`
	Bucket             string `json:"bucket"`
	XAmzAlgorithm      string `json:"X-Amz-Algorithm"`
	XAmzCredential     string `json:"X-Amz-Credential"`
	Policy             string `json:"Policy"`
	XAmzSignature      string `json:"X-Amz-Signature"`
}

type PostIdStruct struct {
	PostId int `json:"postId"`
}

type Filler struct{}
