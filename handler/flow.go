package handler

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	"github.com/softmaxer/localflow/data"
	"github.com/softmaxer/localflow/views"
)

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c, c.Writer)
}

func Router(dbPath string) *gin.Engine {
	db, err := data.InitDB(dbPath, &data.Experiment{}, &data.Competitor{})
	if err != nil {
		log.Fatal("Unable to initialize SQLite Database: ", err.Error())
	}
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		render(ctx, 200, views.Home())
	})
	router.GET("/experiments", func(ctx *gin.Context) {
		render(ctx, 200, views.Experiments())
	})
	router.GET("/experiments/list", func(ctx *gin.Context) {
		getExperiments(ctx, db)
	})
	router.POST("/experiments/add", func(ctx *gin.Context) {
		addExperiment(ctx, db)
	})
	router.GET("/experiments/:id/competitors/list", func(ctx *gin.Context) {
		getCompetitors(ctx, db)
	})
	router.POST("/experiments/:id/competitors/add", func(ctx *gin.Context) {
		addCompetitor(ctx, db)
	})
	router.GET("/experiments/:id", func(ctx *gin.Context) {
		getExperimentById(ctx, db)
	})
	router.GET("/experiments/:id/open", func(ctx *gin.Context) {
		open(ctx, db)
	})
	router.POST("/experiments/:id/run", func(ctx *gin.Context) {
		run(ctx, db)
	})
	router.POST("/experiments/:id/status", func(ctx *gin.Context) {
		status(ctx, db)
	})

	return router
}
