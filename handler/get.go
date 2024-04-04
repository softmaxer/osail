package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
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
	c.IndentedJSON(http.StatusOK, exps)
}
