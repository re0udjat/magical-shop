package main

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *app) readIDParam(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *app) writeJSON(c *gin.Context, status int, data any) {
	if app.config.env == "dev" {
		c.IndentedJSON(status, data)
	} else {
		c.JSON(status, data)
	}
}
