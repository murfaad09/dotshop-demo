package errors

import "fmt"

type Error struct {
	ErrCode    string `json:"error_code,omitempty" example:"ERR000"`
	Err        error  `json:"-"`
	Message    string `json:"message" example:"Something went wrong"`
	StatusCode int    `json:"code" example:"500"`
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Err.Error()
}

func New(msg string, args ...interface{}) *Error {
	return &Error{StatusCode: 500, Message: fmt.Sprintf(msg, args...)}
}

func NewHttp(code int, msg string, args ...interface{}) *Error {
	return &Error{StatusCode: code, Message: fmt.Sprintf(msg, args...)}
}

func Wrap(err error) *Error {
	e := &Error{StatusCode: 500, Err: err}
	if err != nil {
		e.Message = err.Error()
	}
	return e
}

func WrapHttp(code int, err error) *Error {
	e := &Error{StatusCode: code, Err: err}
	if err != nil {
		e.Message = err.Error()
	}
	return e
}

func WrapHttpWithMessage(code int, err error, msg string, args ...interface{}) *Error {
	return &Error{StatusCode: code, Err: err, Message: fmt.Sprintf(msg, args...)}
}

// Internal Server Error with message: something went wrong
func Internal() *Error {
	return New("something went wrong")
}

func (e *Error) WithError(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) WithMessage(msg string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(msg, args...)
	return e
}

func (e *Error) WithCode(code string) *Error {
	e.ErrCode = code
	return e
}

func (e *Error) WithStatus(code int) *Error {
	e.StatusCode = code
	return e
}

func (e Error) Interface() error {
	return &e
}
