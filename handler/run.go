package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
	"github.com/softmaxer/localflow/pkg/board"
	"github.com/softmaxer/localflow/pkg/llm"
)

func run(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var experiment data.Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		db.Model(&experiment).Update("status", "failed")
	}

	db.Model(&experiment).Update("status", "ongoing")
	var competitors []data.Competitor
	db.Where("experiment_id = ?", id).Find(&competitors)

	startPort := 11435
	var models []llm.Model
	for _, competitor := range competitors {
		model := llm.Model{
			Name:   competitor.Name,
			Host:   competitor.Host,
			Port:   startPort,
			Stream: false,
			Rating: 1500,
		}

		models = append(models, model)
		startPort++
	}

	board := board.Board{Competitors: models}
	judge := &llm.Judge{
		Name:   experiment.Judge,
		Host:   "http://localhost",
		Port:   11434,
		Stream: false,
	}
	judge.Init()
	board.Init()
	promptsList := strings.Split(experiment.Prompts, "----")
	log.Printf("Length of prompts found is %d\n", len(promptsList))
	for idx, prompt := range promptsList {
		log.Printf("Processing prompt number %d\n", idx)
		formattedPrompt := fmt.Sprintf("%s\n%s\n", experiment.System, prompt)
		player1, player2 := board.GenerateCompetitors()
		candiateResponse1 := player1.GetCompletion(formattedPrompt)
		candiateResponse2 := player2.GetCompletion(formattedPrompt)
		judgeCompletion := judge.JudgePrompt(candiateResponse1.Response, candiateResponse2.Response)
		var judgement llm.Judgement
		err = llm.ParseJSON(judgeCompletion, &judgement)
		if err != nil {
			db.Model(&experiment).Update("status", "failed")
		}

		player1.UpdateRating(player2, float64(judgement.Result[0]))
		player2.UpdateRating(player1, float64(judgement.Result[1]))
	}

	db.Model(&experiment).Update("status", "finished")
}
