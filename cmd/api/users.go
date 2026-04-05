package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re0udjat/magic-shop/internal/data"
	"github.com/re0udjat/magic-shop/internal/validator"
)

func (app *app) registerUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.badResquestResponse(c, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Set the password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	v := validator.New()

	// Validate the user
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	// Insert the user
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(c, v.Errors)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	// Send the activation email in a background goroutine
	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	app.writeJSON(c, http.StatusCreated, envelope{"user": user})
}
