package structs

// structs used for JSON decoding!
// named after the endpoint that the json stream is coming from

type JsonStruct interface {
	LoggedIn |
		ListEditedProjects |
		ProjectPosts
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
	NumItems int                 `json:"nItems"`
	Items    []ProjectPostsItems `json:"items"`
}

// add rest of stuff
type ProjectPostsItems struct {
	PostId            int           `json:"postId"`
	Headline          string        `json:"headline"`
	PublishedAt       string        `json:"publishedAt"`
	Filename          string        `json:"filename"`
	State             int           `json:"state"`
	NumComments       int           `json:"numComments"`
	NumSharedComments int           `json:"numSharedComments"`
	ContentWarnings   []string      `json:"cws"`
	Tags              []string      `json:"tags"`
	PlainTextBody     string        `json:"plainTextBody"`
	PostingProject    EditedProject `json:"postingProject"`
}
