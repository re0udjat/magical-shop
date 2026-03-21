package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/re0udjat/magic-shop/internal/data"
)

func (app *app) createItemHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Create a new item",
	})
}

func (app *app) getItemHandler(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.writeJSON(c, http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	item := data.Item{
		ID:        id,
		Name:      "Sword of Truth",
		Rarity:    "Legendary",
		CreatedAt: time.Now(),
	}

	app.writeJSON(c, http.StatusOK, item)
}
