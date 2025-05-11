package request

type UpdateSubscriptionPlanRequest struct {
	Name        *string      `json:"name"`
	Price       *int         `json:"price"`
	Feature     *interface{} `json:"feature"`
	DurationDay *int         `json:"duration_day"`
}
