package models

import (
	"testing"
)

var u User = LoginWithCookie("s%3AT2AjSKWV3yE3yPJEyNN2kNZFuzswEib-.zdF3uSoGZfvsSGMvYPZN%2Fo4ea5Efn6oV1zJXAzmBqUw")

func TestUserInfo(t *testing.T) {
	email := u.Email()
	t.Log(email)
	// testing cache timeout
	//time.Sleep(45 * time.Second)
	id := u.Id()
	t.Log(id)
}

// func TestGetEditedProjects(t *testing.T) {
// 	t.Log(u.GetRawEditedProjects())
// }
