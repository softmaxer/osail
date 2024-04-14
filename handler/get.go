package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
	"github.com/softmaxer/localflow/views"
)

func getExperimentById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var experiment data.Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, nil)
	}
	c.IndentedJSON(http.StatusOK, experiment)
}

func getExperiments(c *gin.Context, db *gorm.DB) {
	var exps []data.Experiment
	db.Find(&exps)
	render(c, 200, views.AllExperiments(exps))
}

func getCompetitorById(c *gin.Context, db *gorm.DB) {
	id := c.Param("experiment_id")
	var competitor data.Competitor
	err := db.First(&competitor, "experiment_id = ?", id).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, nil)
	}
	c.IndentedJSON(http.StatusOK, competitor)
}

func getCompetitors(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var competitors []data.Competitor
	err := db.Where("experiment_id = ?", id).Find(&competitors).Error
	if err != nil {
		log.Println("Error retrieving competitors: ", err.Error())
		render(c, 500, views.FailedExpReq())
		return
	}
	render(c, 200, views.ModelsList(competitors, id))
}
