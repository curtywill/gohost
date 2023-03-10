package models

type Project struct {
	Info map[string]interface{}
}

func NewProject(info map[string]interface{}) Project {
	return Project{info}
}

func (p Project) Handle() string {
	return p.Info["handle"].(string)
}
