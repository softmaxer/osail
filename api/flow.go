package api

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c, c.Writer)
}

func Router(dbPath string) *gin.Engine {
	db, err := initDB(dbPath, &Experiment{})
	if err != nil {
		log.Fatal("Unable to initialize SQLite Database: ", err.Error())
	}
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		render(ctx, 200, Home())
	})
	router.GET("/experiments", func(ctx *gin.Context) {
		render(ctx, 200, Experiments())
	})
	router.GET("/experiments/list", func(ctx *gin.Context) {
		getExperiments(ctx, db)
	})
	router.POST("/experiments/add", func(ctx *gin.Context) {
		addExperiment(ctx, db)
	})
	router.GET("/experiments/:id", func(ctx *gin.Context) {
		getExperimentById(ctx, db)
	})
	router.POST("/experiments/:id/run", func(ctx *gin.Context) {
		runExperiment(ctx, db)
	})

	return router
}
