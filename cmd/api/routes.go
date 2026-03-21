package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) routes() http.Handler {
	// Initialize a new gin router instance
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	router.GET("/v1/health", app.healthcheckHandler)

	// API for items
	router.POST("/v1/items", app.createItemHandler)
	router.GET("/v1/items/:id", app.getItemHandler)

	return router
}
