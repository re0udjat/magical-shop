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
	router.Use(app.rateLimit())

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
	router.GET("/v1/items", app.listItemsHandler)
	router.GET("/v1/items/:id", app.showItemHandler)
	router.PATCH("/v1/items/:id", app.updateItemHandler)
	router.DELETE("/v1/items/:id", app.deleteItemHandler)

	// API for users
	router.POST("/v1/users", app.registerUserHandler)

	return router
}
