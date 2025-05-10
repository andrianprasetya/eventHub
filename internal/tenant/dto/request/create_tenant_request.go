package request

type CreateTenantRequest struct {
	Name               string `json:"name" validate:"required"`
	Email              string `json:"email" validate:"required,email,unique=email:tenants"`
	LogoUrl            string `json:"logo_url" validate:"required"`
	Password           string `json:"password" validate:"required,min=8"`
	SubscriptionPlanID string `json:"subscription_plan_id" validate:"required"`
}
