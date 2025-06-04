package errors

import (
	"net/http"

	"weatherApi/internal/common/errors"
)

var (
	ErrCityNotFound        = errors.New(http.StatusNotFound, "City not found", nil)
	ErrInvalidRequest      = errors.New(http.StatusBadRequest, "Invalid Request", nil)
	ErrInternalServerError = errors.New(http.StatusInternalServerError, "Internal server error", nil)
)
