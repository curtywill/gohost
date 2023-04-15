package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gohost/requests"
	"gohost/structs"
	"log"
)

type Project struct {
	Info structs.EditedProject
	u    user
}

func NewProject(info structs.EditedProject, u user) Project {
	return Project{info, u}
}

func (p Project) AvatarShape() string {
	return p.Info.AvatarShape
}

func (p Project) AvatarURL() string {
	return p.Info.AvatarURL
}

func (p Project) HeaderURL() string {
	return p.Info.HeaderURL
}

func (p Project) Pronouns() string {
	return p.Info.Pronouns
}

func (p Project) URL() string {
	return p.Info.URL
}

func (p Project) Privacy() string {
	return p.Info.Privacy
}

func (p Project) Headline() string {
	return p.Info.Dek
}

func (p Project) Description() string {
	return p.Info.Description
}

func (p Project) DisplayName() string {
	return p.Info.DisplayName
}

func (p Project) Handle() string {
	return p.Info.Handle
}

func (p Project) ProjectId() int {
	return p.Info.ProjectId
}

func (p Project) Flags() [][]byte {
	return p.Info.Flags
}

func (p Project) FrequentlyUsedTags() [][]byte {
	return p.Info.FrequentlyUsedTags
}

func (p Project) GetRawPosts(page int) structs.ProjectPosts {
	r := structs.ProjectPosts{}

	data, _ := requests.Fetch("get", fmt.Sprintf("/project/%s/posts?page=%d", p.Handle(), page), p.u.cookie, nil, nil, false)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (p Project) GetPosts(page int) []Post {
	postsRaw := p.GetRawPosts(page)
	posts := []Post{}
	for _, post := range postsRaw.Items {
		posts = append(posts, NewPost(p, post))
	}
	return posts
}

type postRequest struct {
	AdultContent    bool             `json:"adultContent"`
	Blocks          []structs.Blocks `json:"blocks"`
	Tags            []string         `json:"tags"`
	ContentWarnings []string         `json:"cws"`
	Headline        string           `json:"headline"`
	PostState       int              `json:"postState"`
}

func (p Project) Post(adult bool, markdown []Markdown, attachments []Attachment, tags, cws []string, headline string, draft bool) Post {
	markdownLen := len(markdown)
	attachmentLen := len(attachments)

	blocksLen := markdownLen + attachmentLen
	blocks := make([]structs.Blocks, blocksLen)

	for i := range attachments {
		blocks[i] = attachments[i].GetBlock()
	}
	for i := range markdown {
		blocks[i+attachmentLen] = markdown[i].GetBlock()
	}

	postState := 1
	if draft || attachmentLen != 0 {
		postState = 0
	}

	r := postRequest{adult, blocks, tags, cws, headline, postState}
	reqBody, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	jsonHeaders := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	s := structs.PostIdStruct{}
	data, _ := requests.Fetch("POST",
		fmt.Sprintf("/project/%s/posts", p.Handle()),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatal(err)
	}

	if postState == 1 {
		return p.GetPosts(0)[0]
	}

	// draft
	if len(attachments) == 0 {
		return Post{}
	}

	for i := range attachments {
		(&attachments[i]).upload(s.PostId, p)
	}

	for i := range attachments {
		blocks[i] = attachments[i].GetBlock()
	}

	if !draft {
		postState = 1
	}

	r = postRequest{adult, blocks, tags, cws, headline, postState}
	reqBody, err = json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	requests.Fetch("PUT",
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), s.PostId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	if postState == 1 {
		return p.GetPosts(0)[0]
	}

	return Post{}
}
