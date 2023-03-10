package structs

// structs used for JSON decoding!
// named after the endpoint that the json stream is coming from

type JsonStruct interface {
	LoggedIn |
		ListEditedProjects
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
