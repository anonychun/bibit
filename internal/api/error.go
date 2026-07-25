package api

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

type Error struct {
	Status int
	Errors any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func HttpErrorHandler(c *echo.Context, err error) {
	NewResponse(c).SetErrors(err).Send()
}
