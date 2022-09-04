package apperror

type ValidationError map[string][]string

func NewValidationError(err map[string][]string) ValidationError {
	return err
}
func (v ValidationError) Error() string {
	return ValidationErrorName
}
