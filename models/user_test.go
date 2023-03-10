package models

import "testing"

var u User = LoginWithCookie("s%3Au1D0nqXDI3_i3OXt17GwWgrEmwVTMluv.vVaQ2Il3nUM%2FnYA5dcXFqWoLJliURRGrPMmHn0UkQ%2BM")

// func TestLoginWithCookie(t *testing.T) {
// 	u := LoginWithCookie("s%3Au1D0nqXDI3_i3OXt17GwWgrEmwVTMluv.vVaQ2Il3nUM%2FnYA5dcXFqWoLJliURRGrPMmHn0UkQ%2BM")
// 	info := u.userInfo()

// }

func TestUserInfo(t *testing.T) {
	email := u.Email()
	t.Log(email)
}

// func TestGetEditedProjects(t *testing.T) {
// 	u.getEditedProjects()
// }
