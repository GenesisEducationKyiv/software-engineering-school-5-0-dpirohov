package errors

import (
	"net/http"
	"weatherApi/internal/common/errors"
)

var (
	CityNotFoundError   = errors.New(http.StatusNotFound, "City not found", nil)
	InvalidRequestError = errors.New(http.StatusBadRequest, "Invalid Request", nil)
	InternalServerError = errors.New(http.StatusInternalServerError, "Internal server error", nil)
)
