package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/softmaxer/osail/data"
	"github.com/softmaxer/osail/views"
)

func addExperiment(c *gin.Context, db *gorm.DB) {
	var newExp data.ExperimentRequest
	if err := c.Bind(&newExp); err != nil {
		log.Println("The error that occured was: \n", err.Error())
		render(c, http.StatusBadRequest, views.FailedExpReq())
		return
	}

	file, err := c.FormFile("promptpath")
	if err != nil {
		log.Println("Error reading file: ", err.Error())
	}
	promptsFile, err := file.Open()
	if err != nil {
		render(c, http.StatusInternalServerError, views.FailedCreateExp())
	}
	defer promptsFile.Close()
	fileContents, err := io.ReadAll(promptsFile)
	if err != nil {
		render(c, http.StatusInternalServerError, views.FailedCreateExp())
	}

	exp := data.Experiment{
		Name:    newExp.Name,
		Judge:   newExp.Judge,
		System:  newExp.System,
		Prompts: string(fileContents),
		Id:      fmt.Sprintf("exp_%s", uuid.New().String()),
		Status:  "idle",
	}
	err = db.Create(&exp).Error
	if err != nil {
		render(c, http.StatusInternalServerError, views.FailedExpReq())
		return
	}
	render(c, http.StatusCreated, views.PreprendExperiment(exp))
	return
}

func addCompetitor(c *gin.Context, db *gorm.DB) {
	expId := c.Param("id")
	var newCompetitor data.CompetitorRequest
	if err := c.Bind(&newCompetitor); err != nil {
		render(c, http.StatusBadRequest, views.FailedExpReq())
		return
	}
	portNumber, err := strconv.Atoi(newCompetitor.Port)
	if err != nil {
		log.Printf("Error converting port number to integer: %s", err.Error())
	}
	competitor := data.Competitor{
		Name:         newCompetitor.Name,
		Host:         newCompetitor.Host,
		Port:         int32(portNumber),
		Stream:       false,
		Rating:       1500,
		ExperimentId: expId,
	}
	err = db.Create(&competitor).Error
	if err != nil {
		render(c, http.StatusInternalServerError, views.FailedCreateExp())
		return
	}
	render(c, http.StatusCreated, views.PrependModel(competitor))
}
