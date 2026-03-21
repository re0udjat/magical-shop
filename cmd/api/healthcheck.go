package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) healthcheckHandler(c *gin.Context) {

	data := map[string]string{
		"status":  "available",
		"env":     app.config.env,
		"version": version,
	}

	app.writeJSON(c, http.StatusOK, data)
}
