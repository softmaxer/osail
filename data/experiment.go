package data

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Experiment struct {
	gorm.Model
	Name    string
	Id      string
	Judge   string
	System  string
	Prompts string
	Status  string
}

type ExperimentRequest struct {
	Name       string                `form:"name"`
	Judge      string                `form:"judge"`
	System     string                `form:"system"`
	PromptPath *multipart.FileHeader `form:"promptpath"`
}
