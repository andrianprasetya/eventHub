package errors

import "fmt"

type AppError struct {
	Message string
	Code    int
	Err     error
	Expose  bool
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) UnWrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	return e.Code
}

func (e *AppError) ShouldExpose() bool {
	return e.Expose
}

func New(message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
		Expose:  true,
	}
}

func Wrap(err error, message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
		Expose:  false, // default: hidden unless specifically marked
	}
}

// Jika ingin Wrap dan tetap expose (misal errors input)
func WrapExpose(err error, message string, code int) *AppError {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
		Expose:  true,
	}
}
