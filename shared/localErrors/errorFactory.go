package localErrors

import (
	"net/http"
)

func ErrorCustom(code int, message string, meta map[string]any) *ErrorBase {
	return NewError(code, message, meta)
}

func ErrorUnAuthorized() *ErrorBase {
	return NewError(
		http.StatusUnauthorized,
		"UnAuthorized",
		map[string]any{},
	)
}
