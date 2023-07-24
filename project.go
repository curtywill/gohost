package gohost

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/curtywill/gohost/requests"
)

type Project struct {
	editable bool
	info     editedProjectResponse
	u        User
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

func (p Project) getRawPosts(page int) (projectPostsResponse, error) {
	r := projectPostsResponse{}

	data, _, err := requests.Fetch(p.u.client, "get", fmt.Sprintf("/project/%s/posts?page=%d", p.Handle(), page), p.u.cookie, nil, nil, false)
	if err != nil {
		return r, err
	}
	json.Unmarshal(data, &r)

	return r, nil
}

// Returns all the posts on a given page of the current project.
func (p Project) GetPosts(page int) ([]Post, error) {
	postsRaw, err := p.getRawPosts(page)
	if err != nil {
		return nil, err
	}
	posts := make([]Post, len(postsRaw.Items))
	for i, post := range postsRaw.Items {
		posts[i] = Post{p, post}
	}
	return posts, nil
}

type postRequest struct {
	AdultContent    bool             `json:"adultContent"`
	Blocks          []blocksResponse `json:"blocks"`
	Tags            []string         `json:"tags"`
	ContentWarnings []string         `json:"cws"`
	Headline        string           `json:"headline"`
	PostState       int              `json:"postState"`
}

// Makes a post the given project
//
// Returns a Post struct that contains info about the post, including the postId needed for editing!
//
// Returns an empty Post in the case of a draft or an error
func (p Project) Post(adult bool, markdown []Markdown, attachments []Attachment, tags, cws []string, headline string, draft bool) (Post, error) {
	if !p.editable {
		return Post{}, fmt.Errorf("you cannot edit %s", p.Handle())
	}
	// convert nil values to empty slices for marshalling
	if markdown == nil {
		markdown = []Markdown{}
	}
	if attachments == nil {
		attachments = []Attachment{}
	}
	if tags == nil {
		tags = []string{}
	}
	if cws == nil {
		cws = []string{}
	}
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
	reqBody, _ := json.Marshal(r)

	jsonHeaders := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	s := postIdResponse{}
	data, _, err := requests.Fetch(p.u.client, "POST",
		fmt.Sprintf("/project/%s/posts", p.Handle()),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)
	if err != nil {
		return Post{}, nil
	}
	json.Unmarshal(data, &s)

	// go ahead and post if it's not a draft and has no attachments
	if postState == 1 {
		posts, err := p.GetPosts(0)
		return posts[0], err
	}

	// draft
	// TODO: implement drafts
	if len(attachments) == 0 {
		return Post{}, fmt.Errorf("haven't added support for drafts yet")
	}

	for i := range attachments {
		err = (&attachments[i]).upload(p.u.client, s.PostId, p)
		blocks[i] = attachments[i].getBlock()
	}
	if err != nil {
		return Post{}, err
	}

	if !draft {
		postState = 1
	}

	r = postRequest{adult, blocks, tags, cws, headline, postState}
	reqBody, _ = json.Marshal(r)

	requests.Fetch(p.u.client, "PUT",
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), s.PostId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)

	if postState == 1 {
		posts, err := p.GetPosts(0)
		return posts[0], err
	}

	return Post{}, fmt.Errorf("haven't added support for drafts yet")
}

// Edits the post represented by parameter postId
//
// Returns a Post struct that contains info about the post, including the postId needed for editing!
//
// Returns an empty Post in the case of a draft or an error
func (p Project) EditPost(postId int, adult bool, markdown []Markdown, attachments []Attachment, tags, cws []string, headline string, draft bool) (Post, error) {
	if !p.editable {
		return Post{}, fmt.Errorf("you cannot edit %s", p.Handle())
	}
	// convert nil values to empty slices for marshalling
	if markdown == nil {
		markdown = []Markdown{}
	}
	if attachments == nil {
		attachments = []Attachment{}
	}
	if tags == nil {
		tags = []string{}
	}
	if cws == nil {
		cws = []string{}
	}
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
	reqBody, _ := json.Marshal(r)

	jsonHeaders := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	_, _, err := requests.Fetch(p.u.client, "PUT",
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), postId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)
	if err != nil {
		return Post{}, nil
	}
	// go ahead and post if it's not a draft and has no attachments
	if postState == 1 {
		posts, err := p.GetPosts(0)
		return posts[0], err
	}

	// draft
	// TODO: implement drafts
	if len(attachments) == 0 {
		return Post{}, fmt.Errorf("haven't added support for drafts yet")
	}

	for i := range attachments {
		err = (&attachments[i]).upload(p.u.client, postId, p)
		blocks[i] = attachments[i].getBlock()
	}
	if err != nil {
		return Post{}, err
	}

	if !draft {
		postState = 1
	}

	r = postRequest{adult, blocks, tags, cws, headline, postState}
	reqBody, _ = json.Marshal(r)

	_, _, err = requests.Fetch(p.u.client, "PUT",
		fmt.Sprintf("/project/%s/posts/%d", p.Handle(), postId),
		p.u.cookie, jsonHeaders, bytes.NewReader(reqBody), false)
	if err != nil {
		return Post{}, nil
	}
	if postState == 1 {
		posts, err := p.GetPosts(0)
		return posts[0], err
	}

	return Post{}, fmt.Errorf("haven't added support for drafts yet")
}
