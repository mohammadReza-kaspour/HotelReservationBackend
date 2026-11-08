package localErrors

type ErrorBase struct {
	Code    int
	Message string
	Meta    map[string]any
}

func NewError(code int, message string, meta map[string]any) *ErrorBase {
	return &ErrorBase{
		Code:    code,
		Message: message,
		Meta:    meta,
	}
}
func (e *ErrorBase) Error() string {
	return e.Message
}
