package models

import (
	"gohost/structs"
)

type Post struct {
	project Project
	Info    structs.Items
}

func NewPost(p Project, info structs.Items) Post {
	return Post{p, info}
}

func (p Post) Id() int {
	return p.Info.PostId
}
