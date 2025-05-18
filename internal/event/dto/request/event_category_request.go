package request

type EventCategoryPaginateRequest struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `json:"name"`
}
