package main

import (
	"fmt"
	"log"
	"os"

	"github.com/curtywill/gohost"
)

func main() {
	email := os.Getenv("COHOSTEMAIL")
	pass := os.Getenv("COHOSTPASS")

	user, err := gohost.LoginWithPass(nil, email, pass)
	if err != nil {
		log.Fatal(err)
	}
	// each exported struct has getters like this!
	fmt.Println(user.Id())

	project, err := user.DefaultProject()
	if err != nil {
		log.Fatal(err)
	}

	//attachment, err := gohost.AttachmentBlock("eggbug.jpg", "goated mascot")
	if err != nil {
		log.Fatal(err)
	}
	//attachments := []gohost.Attachment{}
	markdown := []gohost.Markdown{gohost.MarkdownBlock("my first gohost post!")}

	post, err := project.Post(false, markdown, nil, []string{"golang", "API", "gohost"}, nil, "cohost is awesome", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post.Id())
}
