package request

type CreateTenantRequest struct {
	Name               string `json:"name" validate:"required"`
	Email              string `json:"email" validate:"required,email,unique=email:tenants"`
	LogoUrl            string `json:"logo_url" validate:"required"`
	Password           string `json:"password" validate:"required,min=8"`
	SubscriptionPlanID string `json:"subscription_plan_id" validate:"required"`
}

type UpdateInformationTenantRequest struct {
	Name    *string `json:"name"`
	LogoUrl *string `json:"logo_url"`
}

type UpdateStatusTenantRequest struct {
	IsActive *int `json:"is_active" validate:"required;numeric"`
}
