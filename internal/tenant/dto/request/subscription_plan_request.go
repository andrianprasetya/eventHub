package request

type CreateSubscriptionPlanRequest struct {
	Name        string      `json:"name" validate:"required"`
	Price       int         `json:"price" validate:"required"`
	Feature     interface{} `json:"feature" validate:"required"`
	DurationDay int         `json:"duration_day" validate:"required"`
}

type UpdateSubscriptionPlanRequest struct {
	Name        *string      `json:"name"`
	Price       *int         `json:"price" validate:"numeric"`
	Feature     *interface{} `json:"feature"`
	DurationDay *int         `json:"duration_day" validate:"numeric""`
}

type SubscriptionPaginateParams struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `query:"name"`
}
