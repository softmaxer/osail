package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
	"github.com/softmaxer/localflow/pkg/board"
	"github.com/softmaxer/localflow/pkg/llm"
	"github.com/softmaxer/localflow/views"
)

func run(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var experiment data.Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		render(ctx, 404, views.ExpProgress("failed"))
	}

	render(ctx, 200, views.ExpProgress("ongoing"))

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
		render(ctx, 500, views.ExpProgress("failed"))
	}

	player1.UpdateRating(player2, float64(judgement.Result[0]))
	player2.UpdateRating(player1, float64(judgement.Result[1]))
	var ratings []float64
	for _, comp := range board.Competitors {
		ratings = append(ratings, comp.Rating)
	}

	render(ctx, 200, views.ExpResults(board.Competitors))
}
