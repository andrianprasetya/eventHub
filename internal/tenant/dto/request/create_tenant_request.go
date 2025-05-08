package request

type CreateTenantRequest struct {
	Name               string `json:"name" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	LogoUrl            string `json:"logo_url" validate:"required"`
	SubscriptionPlanID string `json:"subscription_plan_id" validate:"required"`
}
