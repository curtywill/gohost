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

	// make our first image block with some alt text! can pass an absolute path as well
	attachment, err := gohost.AttachmentBlock("eggbug.jpg", "goated mascot")
	if err != nil {
		log.Fatal(err)
	}
	// build our slices of attachments and markdown
	attachments := []gohost.Attachment{attachment}
	markdown := []gohost.Markdown{gohost.MarkdownBlock("my first gohost post!")}

	tags := []string{"golang", "API", "gohost"}
	headline := "cohost is awesome"
	post, err := project.Post(false, markdown, attachments, tags, nil, headline, false) // returns our post
	if err != nil {
		log.Fatal(err)
	}
	// made our first post onto the default project, now lets edit
	attachment, err = gohost.AttachmentBlock("gopher.png", "another goated mascot")
	if err != nil {
		log.Fatal(err)
	}
	// can append more attachments and markdown to get added to EditPost
	// can also get blocks using post.Blocks()
	attachments = append(attachments, attachment)
	markdown = append(markdown, gohost.MarkdownBlock("edited my first gohost post!"))

	// edit post using the Id() getter from Post
	_, err = project.EditPost(post.Id(), false, markdown, attachments, tags, nil, headline+"!!!", false)
	if err != nil {
		log.Fatal(err)
	}
	// all done! view this post at https://cohost.org/curtywill/post/1378909-cohost-is-awesome
}
