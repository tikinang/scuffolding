package skelet

import (
	"fmt"
	"net/http"
)

type ApiException struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Meta       map[string]any `json:"meta,omitempty"`
}

func (r ApiException) Error() string {
	return fmt.Sprintf("api exception: %d: %s", r.StatusCode, r.Message)
}

func NewApiException(code int, message string, meta map[string]any) ApiException {
	return ApiException{
		StatusCode: code,
		Message:    message,
		Meta:       meta,
	}
}

func NewBadRequestException(reason string) ApiException {
	return ApiException{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad request.",
		Meta: map[string]any{
			"reason": reason,
		},
	}
}

func NewApiExceptionFromError(err error) ApiException {
	return ApiException{
		StatusCode: http.StatusInternalServerError,
		Message:    "Unexpected error occurred.",
		Meta: map[string]any{
			"error_message": err.Error(),
		},
	}
}
