/*

Package error implements a simple way to
create errors to with more context and which
can be displayed to the user upon failure.

*/
package error

import (
	"net/http"
)

// HTTPTypeError defines a struct to return
// errors with some additional
// information
type HTTPTypeError struct {
	StatusCode    int
	StatusMessage string
	Error         string
}

// ErrBadRequest bad request error
func ErrBadRequest(err error) *HTTPTypeError {

	return &HTTPTypeError{
		StatusCode:    http.StatusBadRequest,
		StatusMessage: http.StatusText(http.StatusBadRequest),
		Error:         err.Error(),
	}
}

// ErrInternalServerError bad request error
func ErrInternalServerError(err error) *HTTPTypeError {

	return &HTTPTypeError{
		StatusCode:    http.StatusInternalServerError,
		StatusMessage: http.StatusText(http.StatusInternalServerError),
		Error:         err.Error(),
	}
}
