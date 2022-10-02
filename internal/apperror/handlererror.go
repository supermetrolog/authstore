package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound = NewHandlerErrorWithMessage(errors.New("Entity not found"), "Not found", http.StatusNotFound)
)

const (
	ValidationErrorName = "Validation error"
	HandlerErrorName    = "Handler error"
	LoginErrorName      = "Login error"
	AuthErrorName       = "Login error"
)

type HandlerError struct {
	AppError
	Code       int    `json:"code"`
	StatusText string `json:"status_text"`
	Name       string `json:"name"`
}

func NewHandlerError(err error, code int) *HandlerError {
	return &HandlerError{
		AppError: AppError{
			OriginalError: err,
		},
		Code:       code,
		StatusText: http.StatusText(code),
		Name:       HandlerErrorName,
	}
}
func NewHandlerErrorWithMessage(err error, message string, code int) *HandlerError {
	he := NewHandlerError(err, code)
	he.Message = message
	return he
}
func (he HandlerError) Error() string {
	return HandlerErrorName
}

func (he HandlerError) StatusCode() int {
	return he.Code
}

func (he HandlerError) OriginError() error {
	return he.OriginalError
}

func (he *HandlerError) SetStatusCode(code int) {
	he.Code = code
	he.SetStatusText()
}

func (he *HandlerError) SetStatusText() {
	he.StatusText = http.StatusText(he.Code)
}

func (he *HandlerError) SetName(name string) {
	he.Name = name
}
