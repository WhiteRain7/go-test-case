package apperr

import (
	"fmt"
	"net/http"
)

type CodedError struct {
	code    int
	message string
	err     error
}

func (e *CodedError) Code() int      { return e.code }
func (e *CodedError) Error() string  { return e.message }
func (e *CodedError) Unwrap() error  { return e.err }
func (e *CodedError) String() string { return fmt.Sprintf("%d: %s", e.code, e.message) }

func New(code int, err error, format string, args ...any) *CodedError {
	return &CodedError{code: code, message: fmt.Sprintf(format, args...), err: err}
}

func Unknown(err error) *CodedError {
	return &CodedError{code: http.StatusInternalServerError, message: err.Error(), err: err}
}

func ErrorResponse(err error, res http.ResponseWriter) int {
	if coded, ok := err.(*CodedError); ok {
		http.Error(res, coded.Error(), coded.Code())
		return coded.Code()
	} else {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return http.StatusInternalServerError
	}
}
