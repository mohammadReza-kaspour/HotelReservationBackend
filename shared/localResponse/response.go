package localResponse

import "uuid"

type Response[T any] struct {
	IsSuccess bool
	TraceId   string
	Code      int
	Message   string
	Data      T
}

func NewResponse[T any](code int, message string, payload *T) *Response[T] {
	response := &Response[T]{
		IsSuccess: true,
		TraceId:   uuid.NewV4().String(),
		Code:      code,
		Message:   message,
		Data:      *payload,
	}
	return response
}
