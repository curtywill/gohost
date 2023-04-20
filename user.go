package gohost

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gohost/requests"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	cookie string
	client *http.Client
}

// Start a user session with your connect.sid cookie from Cohost
// If you don't have any special needs for an HTTP client, pass it as nil
func LoginWithCookie(client *http.Client, cookie string) (User, error) {
	if client == nil {
		client = http.DefaultClient
	}
	u := User{cookie, client}
	// if the error returned is nil, we've successfully logged in
	_, err := u.userInfo()
	return u, err
}

// Start a user session with your Cohost email and password
// If you don't have any special needs for an HTTP client, pass it as nil
func LoginWithPass(client *http.Client, email, password string) (User, error) {
	// i know nothing of cryptography so shoutout to @valknight (for code samples) and @iliana (for algorithm)
	// https://cohost.org/iliana/post/180187-eggbug-rs-v0-1-3-d
	if client == nil {
		client = http.DefaultClient
	}
	r := saltResponse{}
	data, _, err := requests.Fetch(client, "GET", fmt.Sprintf("/login/salt?email=%s", email), "", nil, nil, false)
	if err != nil {
		return User{}, err
	}

	json.Unmarshal(data, &r)

	salt := strings.ReplaceAll(r.Salt, "-", "A")
	salt = strings.ReplaceAll(salt, "_", "A")
	salt += "=="

	saltDecoded, _ := base64.StdEncoding.DecodeString(salt)

	hash := pbkdf2.Key([]byte(password), []byte(saltDecoded), 200000, 128, crypto.SHA384.New)
	clientHash := base64.StdEncoding.EncodeToString(hash)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
	}

	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("clientHash", clientHash)

	_, header, err := requests.Fetch(client, "POST", "/login", "", headers, strings.NewReader(formData.Encode()), true)
	if err != nil {
		return User{}, err
	}
	cookieHeader := header.Get("set-cookie")

	s := strings.Split(cookieHeader, "=")[1]
	cookie := strings.Split(s, ";")[0]
	u := User{cookie, client}
	_, err = u.userInfo()
	// if the error returned is nil, we've successfully logged in
	return u, err
}

func (u User) userInfo() (loggedInResponse, error) {
	r := loggedInResponse{}
	data, _, err := requests.FetchTrpc(u.client, "login.loggedIn", u.cookie, nil)
	if err != nil {
		return r, err
	}
	json.Unmarshal(data, &r)

	return r, nil
}

func (u User) Id() int {
	// Gonna ignore the userInfo() errors for the getters, for simplification
	info, _ := u.userInfo()
	return info.ProjectId
}

func (u User) Email() string {
	info, _ := u.userInfo()
	return info.Email
}

func (u User) ProjectId() int {
	info, _ := u.userInfo()
	return info.ProjectId
}

func (u User) ModMode() bool {
	info, _ := u.userInfo()
	return info.ModMode
}

func (u User) ReadOnly() bool {
	info, _ := u.userInfo()
	return info.ReadOnly
}

func (u User) Activated() bool {
	info, _ := u.userInfo()
	return info.Activated
}

// Returns your first project.
func (u User) DefaultProject() (Project, error) {
	projects, err := u.getRawEditedProjects()
	if err != nil {
		return Project{}, nil
	}
	defaultProject := projects.Projects[0]
	return Project{defaultProject, u}, nil
}

// Retrieve one of your projects by its handle
func (u User) GetProject(handle string) (Project, error) {
	projects, err := u.getRawEditedProjects()
	if err != nil {
		return Project{}, err
	}
	for _, project := range projects.Projects {
		if project.Handle == handle {
			return Project{project, u}, nil
		}
	}
	return Project{}, fmt.Errorf("no such project found")
}

func (u User) getRawEditedProjects() (listEditedProjectsResponse, error) {
	r := listEditedProjectsResponse{}
	data, _, err := requests.FetchTrpc(u.client, "projects.listEditedProjects", u.cookie, nil)
	if err != nil {
		return listEditedProjectsResponse{}, err
	}
	json.Unmarshal(data, &r)

	return r, nil
}

// Lists all the projects for an authenticated user
func (u User) GetEditedProjects() ([]Project, error) {
	projectsRaw, err := u.getRawEditedProjects()
	if err != nil {
		return nil, err
	}
	projects := []Project{}
	for _, project := range projectsRaw.Projects {
		projects = append(projects, Project{project, u})
	}

	return projects, nil
}
