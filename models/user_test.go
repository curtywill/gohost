package models

import (
	"testing"
)

var u User = LoginWithCookie("s%3AorZ3lHts7fDPdw7cfKSkAnzBokgRzcPC.npGItUqkzNvgPgEj1UTdlGdo8QJc2fK%2FAOHwJ6KV2T8")

// func TestUserInfo(t *testing.T) {
// 	email := u.Email()
// 	t.Log(email)
// 	// testing cache timeout
// 	//time.Sleep(45 * time.Second)
// 	id := u.Id()
// 	t.Log(id)
// }

// func TestGetEditedProjects(t *testing.T) {
// 	t.Log(u.GetRawEditedProjects())
// }

// func TestGetProjectArray(t *testing.T) {
// 	projects := u.GetEditedProjects()
// 	t.Log(projects)
// 	//t.Log(projects[0].Dek())
// 	//t.Log(projects[1].ProjectId())
// }

// func TestGetRawPosts(t *testing.T) {
// 	projects := u.GetEditedProjects()
// 	t.Log(projects[1].GetRawPosts(0))
// }

func TestPost(t *testing.T) {
	defaultProject := u.GetEditedProjects()[1]
	markdown := []Markdown{NewMarkdown("hello my friends")}
	attachments := []Attachment{NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-03\\chrome_yIz40Aetir.png", "")}

	post := defaultProject.Post(false, markdown, attachments,
		[]string{"golang", "api", "software-development"}, []string{},
		"AUTOMATED POST WITH ATTACHMENT", false)

	t.Log(post.Info.Blocks[0].Attachment)
}
