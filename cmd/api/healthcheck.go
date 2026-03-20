package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) healthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"env":     app.config.env,
		"version": version,
	})
}
