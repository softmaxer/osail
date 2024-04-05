package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
	"github.com/softmaxer/localflow/views"
)

func status(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var experiment data.Experiment
	err := db.First(&experiment, "id = ?", id).Error
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, nil)
	}
	render(c, 200, views.ExpProgress(experiment.Status))
}
