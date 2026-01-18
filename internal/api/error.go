package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Status int
	Errors any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func HttpErrorHandler(err error, c echo.Context) {
	NewResponse(c).SetErrors(err).Send()
}
