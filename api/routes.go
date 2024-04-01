package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/pkg/board"
	"github.com/softmaxer/localflow/pkg/llm"
)

func addExperiment(c *gin.Context, db *gorm.DB) {
	var newExp Experiment
	if err := c.BindJSON(&newExp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
	}

	err := db.Create(&newExp).Error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
	c.IndentedJSON(http.StatusCreated, newExp)
}

func deleteExperiment(c *gin.Context) {
	// TODO
}

func getExperimentById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var experiment Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, nil)
	}
	c.IndentedJSON(http.StatusOK, experiment)
}

func getExperiments(c *gin.Context, db *gorm.DB) {
	var exps []Experiment
	db.Find(&exps)
	c.IndentedJSON(http.StatusOK, exps)
}

func runExperiment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var experiment Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, nil)
	}
	experiment.Status = "ongoing"
	startPort := 11435
	var models []llm.Model
	competitorsList := strings.Split(experiment.Competitors, ",")
	for _, competitorName := range competitorsList {
		model := llm.Model{
			Name:   competitorName,
			Host:   "http://localhost",
			Port:   startPort,
			Stream: false,
			Rating: 1500,
		}

		models = append(models, model)
		startPort++
	}

	board := board.Board{Competitors: models}
	judge := &llm.Judge{Name: "mistral", Host: "http://localhost", Port: 11434, Stream: false}
	judge.Init()
	board.Init()
	player1, player2 := board.GenerateCompetitors()
	candiateResponse1 := player1.GetCompletion(experiment.Prompt)
	candiateResponse2 := player2.GetCompletion(experiment.Prompt)
	judgeCompletion := judge.JudgePrompt(candiateResponse1.Response, candiateResponse2.Response)
	var judgement llm.Judgement
	err = llm.ParseJSON(judgeCompletion, &judgement)
	if err != nil {
		experiment.Status = "failed"
		c.IndentedJSON(http.StatusInternalServerError, experiment)
	}

	player1.UpdateRating(player2, float64(judgement.Result[0]))
	player2.UpdateRating(player1, float64(judgement.Result[1]))

	// board.SortRatings()
	var standings []string
	for _, competitor := range board.Competitors {
		standings = append(standings, competitor.Name)
	}

	db.Model(&experiment).Where("id = ?", id).Update("competitors", strings.Join(standings, ","))
	experiment.Status = "finished"
	c.IndentedJSON(http.StatusOK, experiment)
}
