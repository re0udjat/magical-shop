package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) routes() http.Handler {
	// Initialize a new gin router instance
	router := gin.Default()

	router.GET("/v1/health", app.healthcheckHandler)

	return router
}
