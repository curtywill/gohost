package models

import (
	"gohost/requests"
	"gohost/structs"
	//"log"
)

type User struct {
	cookie string
}

func LoginWithCookie(cookie string) User {
	u := User{cookie}
	// if userInfo doesn't fail, we've successfully logged in
	u.userInfo()
	return u
}

func (u User) userInfo() structs.LoggedIn {
	r := structs.LoggedIn{}
	requests.FetchTrpc("login.loggedIn", u.cookie, &r)

	return r
}

func (u User) Id() int {
	return u.userInfo().UserId
}

func (u User) Email() string {
	return u.userInfo().Email
}

func (u User) ProjectId() int {
	return u.userInfo().ProjectId
}

func (u User) ModMode() bool {
	return u.userInfo().ModMode
}

func (u User) ReadOnly() bool {
	return u.userInfo().ReadOnly
}

func (u User) Activated() bool {
	return u.userInfo().Activated
}

// func (u User) DefaultProject() Project {
// 	return NewProject()
// }

func (u User) GetRawEditedProjects() structs.ListEditedProjects {
	r := structs.ListEditedProjects{}
	requests.FetchTrpc("projects.listEditedProjects", u.cookie, &r)

	return r
}

func (u User) GetEditedProjects() []Project {
	projectsRaw := u.GetRawEditedProjects()
	projects := []Project{}
	for _, project := range projectsRaw.Projects {
		projects = append(projects, NewProject(project, u))
	}

	return projects
}
