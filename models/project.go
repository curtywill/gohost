package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gohost/requests"
	"log"
)

type Project struct {
	info editedProjectResponse
	u    User
}

func (p Project) AvatarShape() string {
	return p.info.AvatarShape
}

func (p Project) AvatarURL() string {
	return p.info.AvatarURL
}

func (p Project) HeaderURL() string {
	return p.info.HeaderURL
}

func (p Project) Pronouns() string {
	return p.info.Pronouns
}

func (p Project) URL() string {
	return p.info.URL
}

func (p Project) Privacy() string {
	return p.info.Privacy
}

func (p Project) Headline() string {
	return p.info.Dek
}

func (p Project) Description() string {
	return p.info.Description
}

func (p Project) DisplayName() string {
	return p.info.DisplayName
}

func (p Project) Handle() string {
	return p.info.Handle
}

func (p Project) ProjectId() int {
	return p.info.ProjectId
}

func (p Project) Flags() [][]byte {
	return p.info.Flags
}

func (p Project) FrequentlyUsedTags() []string {
	return p.info.FrequentlyUsedTags
}

func (p Project) GetRawPosts(page int) projectPostsResponse {
	r := projectPostsResponse{}

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
		posts = append(posts, Post{p, post})
	}
	return posts
}

type postRequest struct {
	AdultContent    bool             `json:"adultContent"`
	Blocks          []blocksResponse `json:"blocks"`
	Tags            []string         `json:"tags"`
	ContentWarnings []string         `json:"cws"`
	Headline        string           `json:"headline"`
	PostState       int              `json:"postState"`
}

func (p Project) Post(adult bool, markdown []Markdown, attachments []Attachment, tags, cws []string, headline string, draft bool) Post {
	markdownLen := len(markdown)
	attachmentLen := len(attachments)

	blocksLen := markdownLen + attachmentLen
	blocks := make([]blocksResponse, blocksLen)

	for i := range attachments {
		blocks[i] = attachments[i].getBlock()
	}
	for i := range markdown {
		blocks[i+attachmentLen] = markdown[i].getBlock()
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

	s := postIdResponse{}
	data, _ := requests.Fetch("POST",
		fmt.Sprintf("/project/%s/posts", p.Handle()),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatal(err)
	}

	// go ahead and post if it's not a draft and has no attachments
	if postState == 1 {
		return p.GetPosts(0)[0]
	}

	// draft
	// TODO: implement drafts
	if len(attachments) == 0 {
		return Post{}
	}

	for i := range attachments {
		(&attachments[i]).upload(s.PostId, p)
		blocks[i] = attachments[i].getBlock()
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

func (p Project) EditPost(postId int, adult bool, markdown []Markdown, attachments []Attachment, tags, cws []string, headline string, draft bool) Post {
	markdownLen := len(markdown)
	attachmentLen := len(attachments)

	blocksLen := markdownLen + attachmentLen
	blocks := make([]blocksResponse, blocksLen)

	for i := range attachments {
		blocks[i] = attachments[i].getBlock()
	}
	for i := range markdown {
		blocks[i+attachmentLen] = markdown[i].getBlock()
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

	requests.Fetch("PUT",
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), postId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	// go ahead and post if it's not a draft and has no attachments
	if postState == 1 {
		return p.GetPosts(0)[0]
	}

	// draft
	// TODO: implement drafts
	if len(attachments) == 0 {
		return Post{}
	}

	for i := range attachments {
		(&attachments[i]).upload(postId, p)
		blocks[i] = attachments[i].getBlock()
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
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), postId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	if postState == 1 {
		return p.GetPosts(0)[0]
	}

	return Post{}
}
