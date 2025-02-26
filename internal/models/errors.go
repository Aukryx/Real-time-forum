package models

type CustomError struct {
	StatusCode int
	Message    string
}

// Error handling function
func NewError(status int, msg string) *CustomError {
	return &CustomError{StatusCode: status, Message: msg}
}
