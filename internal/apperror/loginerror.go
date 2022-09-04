package apperror

import "errors"

type LoginError error

func NewLoginError(message string) LoginError {
	return LoginError(errors.New(message))
}
