package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Response struct {
	echoCtx    *echo.Context
	statusCode int
	body       struct {
		Ok     bool `json:"ok"`
		Meta   any  `json:"meta"`
		Data   any  `json:"data"`
		Errors any  `json:"errors"`
	}
}

func NewResponse(c *echo.Context) *Response {
	return &Response{
		echoCtx:    c,
		statusCode: http.StatusOK,
	}
}

func (r *Response) SetStatus(status int) *Response {
	r.statusCode = status
	return r
}

func (r *Response) SetMeta(meta any) *Response {
	r.body.Meta = meta
	return r
}

func (r *Response) SetData(data any) *Response {
	r.body.Data = data
	return r
}

func (r *Response) SetErrors(err error) *Response {
	switch e := err.(type) {
	case *Error:
		r.SetStatus(e.Status)
		message, ok := e.Errors.(string)
		if ok {
			r.body.Errors = map[string]string{"message": message}
		} else {
			r.body.Errors = e.Errors
		}
	case ValidationError:
		r.SetStatus(http.StatusUnprocessableEntity)
		r.body.Errors = map[string]ValidationError{"params": e}
	case *echo.HTTPError:
		r.SetStatus(e.Code)
		r.body.Errors = map[string]string{"message": e.Message}
	default:
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if statusCode := sc.StatusCode(); statusCode != 0 {
				r.SetStatus(statusCode)
				r.body.Errors = map[string]string{"message": e.Error()}
				break
			}
		}

		r.SetStatus(http.StatusInternalServerError)
		r.body.Errors = map[string]string{"message": "Something went wrong"}
	}

	return r
}

func (r *Response) Send() error {
	r.body.Ok = r.statusCode >= http.StatusOK && r.statusCode < http.StatusMultipleChoices
	if r.body.Ok {
		r.body.Errors = nil
	} else {
		r.body.Data = nil
	}

	return r.echoCtx.JSON(r.statusCode, r.body)
}

func (r *Response) SendMessage(message string) error {
	return r.SetData(map[string]string{"message": message}).Send()
}

func (r *Response) SendOk() error {
	return r.SetStatus(http.StatusOK).SendMessage("ok")
}
