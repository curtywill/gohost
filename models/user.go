package models

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gohost/requests"
	"log"
	"net/url"
	"strings"

	"golang.org/x/crypto/pbkdf2"
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

func LoginWithPass(email, password string) User {
	// i know nothing of cryptography so shoutout to @valknight (for code samples) and @iliana (for algorithm)
	// https://cohost.org/iliana/post/180187-eggbug-rs-v0-1-3-d

	r := saltResponse{}
	data, _ := requests.Fetch("GET", fmt.Sprintf("/login/salt?email=%s", email), "", nil, nil, false)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}

	salt := strings.ReplaceAll(r.Salt, "-", "A")
	salt = strings.ReplaceAll(salt, "_", "A")
	salt += "=="

	saltDecoded, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		log.Fatal(err)
	}

	hash := pbkdf2.Key([]byte(password), []byte(saltDecoded), 200000, 128, crypto.SHA384.New)
	clientHash := base64.StdEncoding.EncodeToString(hash)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
	}

	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("clientHash", clientHash)

	_, header := requests.Fetch("POST", "/login", "", headers, strings.NewReader(formData.Encode()), true)
	cookie := header.Get("set-cookie")

	s := strings.Split(cookie, "=")[1]
	u := User{strings.Split(s, ";")[0]}
	u.userInfo()

	return u
}

func (u User) userInfo() loggedInResponse {
	r := loggedInResponse{}
	data, _ := requests.FetchTrpc("login.loggedIn", u.cookie, nil)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}
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

func (u User) DefaultProject() Project {
	projects := u.getRawEditedProjects()
	defaultProject := projects.Projects[0]
	return Project{defaultProject, u}
}

// retrieve one of your projects by its handle
func (u User) GetProject(handle string) Project {
	projects := u.getRawEditedProjects()
	for _, project := range projects.Projects {
		if project.Handle == handle {
			return Project{project, u}
		}
	}
	return Project{}
}

func (u User) getRawEditedProjects() listEditedProjectsResponse {
	r := listEditedProjectsResponse{}
	data, _ := requests.FetchTrpc("projects.listEditedProjects", u.cookie, nil)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

// lists all the projects for an authenticated user
func (u User) GetEditedProjects() []Project {
	projectsRaw := u.getRawEditedProjects()
	projects := []Project{}
	for _, project := range projectsRaw.Projects {
		projects = append(projects, Project{project, u})
	}

	return projects
}
