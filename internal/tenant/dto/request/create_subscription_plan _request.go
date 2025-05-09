package request

type CreateSubscriptionPlanRequest struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Feature     string `json:"feature"`
	DurationDay int    `json:"duration_day"`
}
