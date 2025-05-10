package request

type CreateSubscriptionPlanRequest struct {
	Name        string      `json:"name" validate:"required"`
	Price       int         `json:"price" validate:"required"`
	Feature     interface{} `json:"feature" validate:"required"`
	DurationDay int         `json:"duration_day" validate:"required"`
}
