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

func (app *app) listItemsHandler(c *gin.Context) {
	var input struct {
		Name   string
		Rarity string
		data.Filters
	}

	v := validator.New()

	qs := c.Request.URL.Query()

	// Extract data from query string, falling back to  defaults of an empty string and
	// an empty slice respectively if they are not provided by the client
	input.Name = app.readString(qs, "name", "")
	input.Rarity = app.readString(qs, "rarity", "")

	// Get the page and page_size query string values as integers
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it's not provided by
	// the client
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "rarity", "price", "-id", "-name", "-rarity", "-price"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	items, metadata, err := app.models.Items.GetAll(input.Name, input.Rarity, input.Filters)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"items": items, "metadata": metadata})
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
