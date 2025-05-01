package response

type APIResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func SuccessResponse(msg string) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Status:  "success",
		Message: msg,
	}
}

func SuccessWithDataResponse[T any](msg string, data T) *APIResponse[T] {
	return &APIResponse[T]{
		Status:  "success",
		Message: msg,
		Data:    data,
	}
}

func ValidationResponse(errs any) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Status:  "validation error",
		Message: "Invalid request",
		Errors:  errs,
	}
}

func ErrorResponse(msg string, errs any) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Status:  "error",
		Message: msg,
		Errors:  errs,
	}
}
