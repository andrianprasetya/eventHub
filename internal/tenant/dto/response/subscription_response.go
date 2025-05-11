package response

type SubscriptionPlanResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Feature     interface{} `json:"feature"`
	Price       int         `json:"price"`
	DurationDay int         `json:"duration_day"`
}

type SubscriptionPlanListItemResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Feature     interface{} `json:"feature"`
	Price       int         `json:"price"`
	DurationDay int         `json:"duration_day"`
}
