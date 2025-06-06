package response

type APIResponse[T any] struct {
	Meta     *Meta     `json:"meta"`
	Data     *Data[T]  `json:"data,omitempty,dive"`
	Errors   any       `json:"errors,omitempty"`
	PageInfo *PageInfo `json:"page_info,omitempty"`
}

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Data[T any] struct {
	Item  T   `json:"item"`
	Items []T `json:"items"`
}

type PageInfo struct {
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
}

func SuccessResponse(code int, msg string) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Meta: &Meta{
			Code:    code,
			Message: msg,
		},
		PageInfo: nil,
	}
}

func SuccessWithDataResponse[T any](code int, msg string, data T) *APIResponse[T] {
	return &APIResponse[T]{
		Meta: &Meta{
			Code:    code,
			Message: msg,
		},
		Data: &Data[T]{
			Item:  data,
			Items: make([]T, 0),
		},
		PageInfo: nil,
	}
}

func SuccessWithPaginateDataResponse[T any](code int, msg string, data []T, page, pageSize int, totalItems int64) *APIResponse[T] {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	var empty T

	return &APIResponse[T]{
		Meta: &Meta{
			Code:    code,
			Message: msg,
		},
		PageInfo: &PageInfo{
			TotalPages: totalPages,
			TotalItems: totalItems,
			Page:       page,
			PageSize:   pageSize,
		},
		Data: &Data[T]{
			Item:  empty,
			Items: data,
		},
	}
}

func ValidationResponse(code int, errs any) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Meta: &Meta{
			Code:    code,
			Message: "Invalid request",
		},
		Errors:   errs,
		PageInfo: nil,
	}
}

func ErrorResponse(code int, msg string, errs any) *APIResponse[struct{}] {
	return &APIResponse[struct{}]{
		Meta: &Meta{
			Code:    code,
			Message: msg,
		},
		Errors:   errs,
		PageInfo: nil,
	}
}
