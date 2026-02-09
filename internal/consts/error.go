package consts

import (
	"net/http"

	"github.com/anonychun/bibit/internal/api"
)

var (
	ErrUnauthorized                  = &api.Error{Status: http.StatusUnauthorized, Errors: "You are not allowed to perform this action"}
	ErrInvalidCredentials            = &api.Error{Status: http.StatusUnauthorized, Errors: "Invalid email or password"}
	ErrEmailAddressAlreadyRegistered = &api.Error{Status: http.StatusConflict, Errors: "Email address already registered"}

	ErrAdminNotFound = &api.Error{Status: http.StatusNotFound, Errors: "Admin not found"}
)
