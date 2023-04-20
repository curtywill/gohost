package gohost

import "time"

// TODO: add rest of getter functions
type Post struct {
	project Project
	info    itemsResponse
}

func (p Post) Id() int {
	return p.info.PostId
}

func (p Post) Headline() string {
	return p.info.Headline
}

func (p Post) PublishedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, p.info.PublishedAt)
	return t
}

func (p Post) Filename() string {
	return p.info.Filename
}

func (p Post) State() int {
	return p.info.State
}

func (p Post) NumComments() int {
	return p.info.NumComments
}

func (p Post) NumSharedComments() int {
	return p.info.NumSharedComments
}

func (p Post) ContentWarnings() []string {
	return p.info.ContentWarnings
}

func (p Post) Tags() []string {
	return p.info.Tags
}

func (p Post) PlainTextBody() string {
	return p.info.PlainTextBody
}

func (p Post) PostingProject(u User) Project {
	return Project{p.info.PostingProject, u}
}

func (p Post) Url() string {
	return p.info.Url
}

func (p Post) Blocks() ([]Markdown, []Attachment) {
	attachments := []Attachment{}
	markdown := []Markdown{}
	for _, block := range p.info.Blocks {
		if block.Type == "attachment" {
			attachments = append(attachments, Attachment{block: *block.Attachment})
		} else {
			markdown = append(markdown, Markdown{*block.Markdown})
		}
	}
	return markdown, attachments
}
