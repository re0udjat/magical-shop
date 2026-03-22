package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/re0udjat/magic-shop/internal/data"
	"github.com/re0udjat/magic-shop/internal/validator"
)

func (app *app) createItemHandler(c *gin.Context) {
	var input struct {
		Name   string        `json:"name"`
		Rarity string        `json:"rarity"`
		Price  data.Currency `json:"price"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.badResquestResponse(c, err)
		return
	}

	// Make a copy of item
	item := &data.Item{
		Name:   input.Name,
		Rarity: input.Rarity,
		Price:  input.Price,
	}

	// Initialize a new Validator
	v := validator.New()

	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
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
