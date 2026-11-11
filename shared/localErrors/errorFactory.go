package localErrors

import (
	"net/http"
)

func NewErrorCustom(code int, message string, meta map[string]any) *ErrorBase {
	return newError(code, message, meta)
}

func NewErrorUnAuthorized(message ...string) *ErrorBase {
	initMessage := "UnAuthorized"
	if len(message) > 0 && message[0] != "" {
		initMessage = message[0]
	}
	return newError(
		http.StatusUnauthorized,
		initMessage,
		map[string]any{},
	)
}
func NewErrorBadRequest(message string) *ErrorBase {
	return newError(
		http.StatusBadRequest,
		message,
		map[string]any{},
	)
}
func NewErrorNotFound(message string) *ErrorBase {
	return newError(
		http.StatusNotFound,
		message,
		map[string]any{},
	)
}
