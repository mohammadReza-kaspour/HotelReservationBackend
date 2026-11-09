package localErrors

import (
	"net/http"
)

func NewErrorCustom(code int, message string, meta map[string]any) *ErrorBase {
	return newError(code, message, meta)
}

func NewErrorUnAuthorized() *ErrorBase {
	return newError(
		http.StatusUnauthorized,
		"UnAuthorized",
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
