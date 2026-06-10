package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Status int
	Errors any
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(httpStatusToGrpcCode(e.Status), e.Error())
}

func httpStatusToGrpcCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case 499: // Client Closed Request
		return codes.Canceled
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		if httpStatus >= 200 && httpStatus < 300 {
			return codes.OK
		}
		return codes.Internal
	}
}

func HttpErrorHandler(c *echo.Context, err error) {
	NewResponse(c).SetErrors(err).Send()
}
