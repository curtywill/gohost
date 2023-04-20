package test

import (
	"gohost/models"
	"testing"
)

// func TestLoginWithCookie(t *testing.T) {
// 	_, err := models.LoginWithCookie(nil, "s:nodcGttIO04V5ISstqG9fLP2CZ-S-cQ7.Zwcsnon0McmckLSQE83R8FdF1NddNuODHsutGeKoFwg")
// 	if err != nil {
// 		t.Log(err)
// 		t.Fatal("Login failed!")
// 	}
// }

func TestLoginWithEmailAndPass(t *testing.T) {
	// email, _ := os.LookupEnv("COHOSTEMAIL")
	// pass, _ := os.LookupEnv("COHOSTPASSWORD")
	// t.Log(email, pass)
	user, err := models.LoginWithPass(nil, "curtisswilliams@rocketmail.com", "Rockyasuho1!")
	if err != nil {
		t.Fatal(err)
	}
	project, err := user.DefaultProject()
	if err != nil {
		t.Fatal(err)
	}
	posts, err := project.GetPosts(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(posts[0].PublishedAt())
}
