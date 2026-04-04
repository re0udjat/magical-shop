package main

import (
	"errors"
	"fmt"
	"net/http"

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

	// Create a record in db and update the item struct with the system-generated info
	err = app.models.Items.Insert(item)
	if err != nil {
		app.errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Add the Location header to let the client know which URL they can find the newly-created resource at
	c.Header("Location", fmt.Sprintf("/v1/items/%d", item.ID))
	app.writeJSON(c, http.StatusCreated, envelope{"item": item})
}

func (app *app) showItemHandler(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	item, err := app.models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"item": item})
}

func (app *app) updateItemHandler(c *gin.Context) {
	// Read the ID from the URL parameters
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	// Get the existing item from the database
	item, err := app.models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	// Read the JSON body
	var input struct {
		Name   *string        `json:"name"`
		Rarity *string        `json:"rarity"`
		Price  *data.Currency `json:"price"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.badResquestResponse(c, err)
		return
	}

	// Update the fields
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Rarity != nil {
		item.Rarity = *input.Rarity
	}
	if input.Price != nil {
		item.Price = *input.Price
	}

	// Validate the updated item
	v := validator.New()
	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	// Save the updated item
	err = app.models.Items.Update(item)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"item": item})
}

func (app *app) deleteItemHandler(c *gin.Context) {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return
	}

	err = app.models.Items.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"message": "item successfully deleted"})
}
