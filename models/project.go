package models

import (
	"fmt"
	"gohost/requests"
	"gohost/structs"
)

type Project struct {
	Info structs.EditedProject
	u    User
}

func NewProject(info structs.EditedProject, u User) Project {
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

	requests.Fetch("get", fmt.Sprintf("/project/%s/posts?page=%d", p.Handle(), page), p.u.cookie, "", false, &r)

	return r
}
