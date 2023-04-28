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

	// default project is the first project you created
	project, err := user.DefaultProject()
	if err != nil {
		log.Fatal(err)
	}

	// every cohost post is made using blocks of attachments and markdown
	// can make both by using composite literals like these!
	attachments := []gohost.Attachment{gohost.Attachment{Filepath: "eggbug.jpg", AltText: "goated mascot"}}
	markdown := []gohost.Markdown{gohost.Markdown{"my first gohost post!"}}

	tags := []string{"golang", "API", "gohost"}
	headline := "cohost is awesome"
	post, err := project.Post(false, markdown, attachments, tags, nil, headline, false) // returns our post
	if err != nil {
		log.Fatal(err)
	}

	// made our first post onto the default project, now lets edit
	// can append more attachments and markdown to get added to EditPost
	// can also get blocks using post.Blocks()
	attachments = append(attachments, gohost.Attachment{Filepath: "gopher.png", AltText: "another goated mascot"})
	markdown = append(markdown, gohost.Markdown{"edited my first gohost post!"})

	// edit post using the Id() getter from Post
	_, err = project.EditPost(post.Id(), false, markdown, attachments, tags, nil, headline+"!!!", false)
	if err != nil {
		log.Fatal(err)
	}
	// all done! view this post at https://cohost.org/curtywill/post/1378909-cohost-is-awesome
}
