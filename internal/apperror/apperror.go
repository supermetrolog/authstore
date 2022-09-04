package apperror

type ResponseError struct {
	Error error `json:"error"`
}
type AppError struct {
	OriginalError error  `json:"original_error"`
	Message       string `json:"message,omitempty"`
}
