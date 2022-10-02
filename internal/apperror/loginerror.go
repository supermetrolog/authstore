package apperror

type LoginError struct {
	Message string `json:"message"`
}

func NewLoginError(message string) LoginError {
	return LoginError{
		Message: message,
	}
}

func (l LoginError) Error() string {
	return l.Message
}
