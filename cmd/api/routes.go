package main

import (
	"expvar"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *app) routes() http.Handler {
	// Initialize a new gin router instance
	router := gin.New()
	router.HandleMethodNotAllowed = true

	// Add middleware
	router.Use(gin.Logger()).
		Use(app.metrics()).
		Use(app.recover()).
		Use(app.enableCORS()).
		Use(app.rateLimit()).
		Use(app.authenticate())

	// Customize error response for router
	router.NoRoute(func(c *gin.Context) {
		app.notFoundResponse(c)
	})
	router.NoMethod(func(c *gin.Context) {
		app.methodNotAllowedResponse(c)
	})

	// Group router
	v1 := router.Group("/v1")

	// Health check endpoint
	v1.GET("/health", app.healthcheckHandler)

	// API for items
	items := v1.Group("/items")
	items.POST("/", app.requirePermission("items:write", app.createItemHandler))
	items.GET("/", app.requirePermission("items:read", app.listItemsHandler))
	items.GET("/:id", app.requirePermission("items:read", app.showItemHandler))
	items.PATCH("/:id", app.requirePermission("items:write", app.updateItemHandler))
	items.DELETE("/:id", app.requirePermission("items:write", app.deleteItemHandler))

	// API for users
	router.POST("/v1/users", app.registerUserHandler)
	router.PUT("/v1/users/activated", app.activateUserHandler)

	// API for tokens
	router.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// API for metrics
	router.GET("/debug/vars", gin.WrapH(expvar.Handler()))

	return router
}
