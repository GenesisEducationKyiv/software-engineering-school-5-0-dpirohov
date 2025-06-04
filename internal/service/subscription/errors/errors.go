package errors

import (
	"net/http"

	"weatherApi/internal/common/errors"
)

var (
	ErrInvalidInput        = errors.New(http.StatusBadRequest, "Invalid input", nil)
	ErrAlreadySubscribed   = errors.New(http.StatusConflict, "Email already subscribed", nil)
	ErrInternalServerError = errors.New(http.StatusInternalServerError, "Internal server error", nil)
	ErrTokenNotFound       = errors.New(http.StatusNotFound, "Token not found", nil)
	ErrInvalidToken        = errors.New(http.StatusBadRequest, "Invalid token", nil)
)
