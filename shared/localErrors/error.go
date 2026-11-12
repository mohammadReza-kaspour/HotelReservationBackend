package localErrors

import "uuid"

type ErrorBase struct {
	IsSuccess bool
	TraceId   string
	Code      int
	Message   string
	Data      map[string]any
}

func newError(code int, message string, data map[string]any) *ErrorBase {
	return &ErrorBase{
		IsSuccess: false,
		TraceId:   uuid.NewV4().String(),
		Code:      code,
		Message:   message,
		Data:      data,
	}
}
func (e *ErrorBase) Error() string {
	return e.Message
}
