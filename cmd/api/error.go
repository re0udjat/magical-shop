package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// A generic helper for logging an error msg along with the current request
// method and URL as attributes in the log entry
func (app *app) logError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// A generic helper for sending JSON-formatted error msg to the client with a given status code
func (app *app) errorResponse(c *gin.Context, status int, msg any) {
	env := envelope{"error": msg}
	app.writeJSON(c, status, env)
}

// Used when app encounters an unexpected problem at runtime
// Logs the detailed error msg
// Sends a 500 Internal Server Error status code and JSON response to client
func (app *app) serverErrorResponse(c *gin.Context, err error) {
	app.logError(c, err)
	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(c, http.StatusInternalServerError, msg)
}

// Send a 400 Bad Request status code and JSON response to client
func (app *app) badResquestResponse(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}

// Send a 404 Not Found status code and JSON response to client
func (app *app) notFoundResponse(c *gin.Context) {
	msg := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, msg)
}

// Send a 405 Method Not Allowed status code and JSON response to client
func (app *app) methodNotAllowedResponse(c *gin.Context) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	app.errorResponse(c, http.StatusMethodNotAllowed, msg)
}

// Send a 409 Conflict status code and JSON response to client
func (app *app) editConflictResponse(c *gin.Context) {
	msg := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(c, http.StatusConflict, msg)
}

// Send a 422 Unprocessable Entity status code and JSON response to client
func (app *app) failedValidationResponse(c *gin.Context, errors map[string]string) {
	app.errorResponse(c, http.StatusUnprocessableEntity, errors)
}

// Send a 429 Too Many Requests status code and JSON response to client
func (app *app) rateLimitExceededResponse(c *gin.Context) {
	msg := "rate limit exceeded"
	app.errorResponse(c, http.StatusTooManyRequests, msg)
}

// Send a 401 Unauthorized status code and JSON response to client
func (app *app) invalidCredentialsResponse(c *gin.Context) {
	msg := "invalid authentication credentials"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}

func (app *app) invalidAuthenticationTokenResponse(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	msg := "invalid or missing authentication token"
	app.errorResponse(c, http.StatusUnauthorized, msg)
}
