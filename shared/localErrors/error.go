package localErrors

import "uuid"

type ErrorBase struct {
	IsSuccess bool
	TraceId   string
	Code      int
	Message   string
	Meta      map[string]any
}

func newError(code int, message string, meta map[string]any) *ErrorBase {
	return &ErrorBase{
		IsSuccess: false,
		TraceId:   uuid.NewV4().String(),
		Code:      code,
		Message:   message,
		Meta:      meta,
	}
}
func (e *ErrorBase) Error() string {
	return e.Message
}
