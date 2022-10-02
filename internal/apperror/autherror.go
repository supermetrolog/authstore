package apperror

type AuthError struct {
	Message string `json:"message"`
}

func NewAuthError(message string) AuthError {
	return AuthError{
		Message: message,
	}
}

func (a AuthError) Error() string {
	return a.Message
}
