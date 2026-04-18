package apperr

type Code int

const (
	CodeInvalidRequest Code = iota + 1
	CodeUnauthorized
	CodeNotFound
	CodeConflict
	CodeInternal
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Message
	}
	if e.Message == "" {
		return e.Err.Error()
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(code Code, message string) error {
	return &Error{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) error {
	return &Error{Code: code, Message: message, Err: err}
}

func InvalidRequest(message string) error { return New(CodeInvalidRequest, message) }
func Unauthorized(message string) error   { return New(CodeUnauthorized, message) }
func NotFound(message string) error       { return New(CodeNotFound, message) }
func Conflict(message string) error       { return New(CodeConflict, message) }
func Internal(err error) error            { return Wrap(CodeInternal, "internal server error", err) }
