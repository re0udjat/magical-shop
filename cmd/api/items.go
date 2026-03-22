package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/re0udjat/magic-shop/internal/data"
)

func (app *app) createItemHandler(c *gin.Context) {
	var input struct {
		Name   string `json:"name"`
		Rarity string `json:"rarity"`
		Price  int64  `json:"price"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		app.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{
		"msg": fmt.Sprintf("%+v", input),
	})
}

func (app *app) getItemHandler(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item := data.Item{
		ID:        id,
		Name:      "Sword of Truth",
		Rarity:    "Legendary",
		CreatedAt: time.Now(),
	}

	app.writeJSON(c, http.StatusOK, envelope{"item": item})
}
