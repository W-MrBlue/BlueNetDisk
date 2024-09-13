package e

import "net/http"

const (
	Success             = http.StatusOK
	InvalidParams       = http.StatusBadRequest
	InternalServerError = http.StatusInternalServerError
)
