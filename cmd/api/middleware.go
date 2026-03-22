package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (app *app) recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a deferred function (which will always be run in the event of a panic as Go unwind the stack)
		defer func() {
			// Builtin recover function to check if there has been a panic or not
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response to make Go's HTTP server
				// automatically close the current connection after a response has been sent
				c.Header("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%s", err))
			}
		}()
		c.Next()
	}
}
