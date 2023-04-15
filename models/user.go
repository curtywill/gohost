package models

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	//"fmt"
	"gohost/requests"
	"gohost/structs"
	"log"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	//"log"
)

type user struct {
	cookie string
}

func LoginWithCookie(cookie string) user {
	u := user{cookie}
	// if userInfo doesn't fail, we've successfully logged in
	u.userInfo()
	return u
}

func LoginWithPass(email, password string) user {
	// i know nothing of cryptography so shoutout to @valknight (for code samples) and @iliana (for algorithm)
	// https://cohost.org/iliana/post/180187-eggbug-rs-v0-1-3-d

	r := structs.SaltResponse{}
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
	u := user{strings.Split(s, ";")[0]}
	u.userInfo()

	return u
}

func (u user) userInfo() structs.LoggedIn {
	r := structs.LoggedIn{}
	data, _ := requests.FetchTrpc("login.loggedIn", u.cookie, nil)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func (u user) Id() int {
	return u.userInfo().UserId
}

func (u user) Email() string {
	return u.userInfo().Email
}

func (u user) ProjectId() int {
	return u.userInfo().ProjectId
}

func (u user) ModMode() bool {
	return u.userInfo().ModMode
}

func (u user) ReadOnly() bool {
	return u.userInfo().ReadOnly
}

func (u user) Activated() bool {
	return u.userInfo().Activated
}

// func (u User) DefaultProject() Project {
// 	return NewProject()
// }

func (u user) GetRawEditedProjects() structs.ListEditedProjects {
	r := structs.ListEditedProjects{}
	data, _ := requests.FetchTrpc("projects.listEditedProjects", u.cookie, nil)

	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func (u user) GetEditedProjects() []Project {
	projectsRaw := u.GetRawEditedProjects()
	projects := []Project{}
	for _, project := range projectsRaw.Projects {
		projects = append(projects, NewProject(project, u))
	}

	return projects
}
