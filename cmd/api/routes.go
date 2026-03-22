package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) routes() http.Handler {
	// Initialize a new gin router instance
	router := gin.New()
	router.HandleMethodNotAllowed = true

	// Add middleware
	router.Use(gin.Logger())
	router.Use(app.recover())

	// Customize error response for router
	router.NoRoute(func(c *gin.Context) {
		app.notFoundResponse(c)
	})
	router.NoMethod(func(c *gin.Context) {
		app.methodNotAllowedResponse(c)
	})

	// Health check endpoint
	router.GET("/v1/health", app.healthcheckHandler)

	// API for items
	router.POST("/v1/items", app.createItemHandler)
	router.GET("/v1/items/:id", app.getItemHandler)

	return router
}
