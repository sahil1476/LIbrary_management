package handlers

import (
	"main/connection"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var lib models.Library

func Createlibrary(c *gin.Context) {
	if err := c.ShouldBindJSON(&lib); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := connection.DB.Create(&lib).Error

	if err != nil {
		c.JSON(400, "Library not created in DATA-BASE")
		return
	}
	c.JSON(http.StatusOK, "Library Created")
}
