package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/osail/data"
	"github.com/softmaxer/osail/pkg/board"
	"github.com/softmaxer/osail/pkg/llm"
	"github.com/softmaxer/osail/views"
)

func updateStatus(ctx *gin.Context) {
	render(ctx, 200, views.ExpProgress("ongoing"))
}

func run(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var experiment data.Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		db.Model(&experiment).Update("status", "failed")
	}

	db.Model(&experiment).Update("status", "ongoing")
	var competitors []data.Competitor
	db.Model(&data.Competitor{}).Where("experiment_id = ?", id).Find(&competitors)

	board := board.Board{Competitors: competitors}
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
		db.Model(player1).Where("name = ?", player1.Name).Update("rating", player1.Rating)
		player2.UpdateRating(player1, float64(judgement.Result[1]))
		db.Model(player2).Where("name = ?", player2.Name).Update("rating", player2.Rating)
	}

	db.Model(&experiment).Update("status", "finished")
	render(ctx, 200, views.HiddenUpdate(id))
}
