package errors

import "net/http"

var (
	ErrInvalidCredentials = New("invalid email or password", http.StatusBadRequest)
	ErrInternalServer     = New("internal server errors", http.StatusInternalServerError)
)
