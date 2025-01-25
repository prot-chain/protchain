package function

import (
	"net/http"
	"protchain/internal/value"
)

func StatusCode(status string) int {
	switch status {
	case value.Success:
		return http.StatusOK
	case value.NotFound:
		return http.StatusNotFound
	case value.Created:
		return http.StatusCreated
	case value.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
