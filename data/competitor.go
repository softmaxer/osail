package data

import "gorm.io/gorm"

type Competitor struct {
	gorm.Model
	Name         string
	Host         string
	Port         int32
	Stream       bool
	Rating       float64
	ExperimentId string
}

type CompetitorRequest struct {
	Name string `form:"name"`
	Host string `form:"host"`
	Port string `form:"port"`
}
