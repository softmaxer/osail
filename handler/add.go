package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/softmaxer/localflow/data"
)

func addExperiment(c *gin.Context, db *gorm.DB) {
	var newExp data.Experiment
	if err := c.BindJSON(&newExp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
	}

	err := db.Create(&newExp).Error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
	c.IndentedJSON(http.StatusCreated, newExp)
}
