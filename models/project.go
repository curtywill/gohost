package models

import "gohost/structs"

type Project struct {
	Info structs.EditedProject
}

func NewProject(info structs.EditedProject) Project {
	return Project{info}
}

func (p Project) Handle() string {
	return p.Info.Handle
}
