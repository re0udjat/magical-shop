package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *app) healthcheckHandler(c *gin.Context) {

	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"env":     app.config.env,
			"version": version,
		},
	}

	time.Sleep(10 * time.Second)

	app.writeJSON(c, http.StatusOK, data)
}
