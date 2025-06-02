package errors

import (
	"net/http"
	"weatherApi/internal/common/errors"
)

var (
	InvalidInputError   = errors.New(http.StatusBadRequest, "Invalid input", nil)
	AlreadySubscribed   = errors.New(http.StatusConflict, "Email already subscribed", nil)
	InternalServerError = errors.New(http.StatusInternalServerError, "Internal server error", nil)
	TokenNotFound       = errors.New(http.StatusNotFound, "Token not found", nil)
	InvalidToken        = errors.New(http.StatusBadRequest, "Invalid token", nil)
)
