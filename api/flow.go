package api

import (
	"log"
	"net/http"

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
	router.GET("/experiments/:id/open", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var experiment Experiment
		err := db.First(&experiment, "id = ?", id).Error
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, nil)
		}

		render(ctx, 200, ShowExperiment(experiment))
	})
	router.POST("/experiments/:id/run", func(ctx *gin.Context) {
		runExperiment(ctx, db)
	})

	return router
}
