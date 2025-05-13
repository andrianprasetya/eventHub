package request

type UpdateSubscriptionPlanRequest struct {
	Name        *string      `json:"name"`
	Price       *int         `json:"price" validate:"numeric"`
	Feature     *interface{} `json:"feature"`
	DurationDay *int         `json:"duration_day" validate:"numeric""`
}
