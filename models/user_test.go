package models

import "testing"

// var u user = LoginWithCookie("s%3AorZ3lHts7fDPdw7cfKSkAnzBokgRzcPC.npGItUqkzNvgPgEj1UTdlGdo8QJc2fK%2FAOHwJ6KV2T8")

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

// func TestPost(t *testing.T) {
// 	defaultProject := u.GetEditedProjects()[1]

// 	attachments := make([]Attachment, 4)

// 	attachments[0] = NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-03\\chrome_5cAlEfQ29S.png", "tsurgi serving")
// 	attachments[1] = NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-03\\chrome_pv179u36J4.png", "rare joshu serve")
// 	attachments[2] = NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-03\\chrome_VR12RGfp7E.png", "bae")
// 	attachments[3] = NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-03\\chrome_vdhusFhtVA.png", "angels")

// 	defaultProject.Post(false, nil, attachments,
// 		[]string{"golang", "api", "software-development"}, []string{},
// 		"jojo dump", false)
// }

func TestLoginWithPass(t *testing.T) {
	u := LoginWithPass("email", "password")
	project := u.GetProject("curtypage")
	post := project.GetPosts(0)[2]

	markdown, attachments := post.Blocks()
	attachments = append(attachments, NewAttachment("C:\\Users\\curty\\Documents\\ShareX\\Screenshots\\2023-02\\chrome_hcLD39Nquc.png", "i added this in post"))
	editedPost := project.EditPost(post.Id(), false, markdown, attachments,
		[]string{"golang", "api", "software development"}, []string{}, post.Headline(), false)
	t.Log(editedPost)
}
