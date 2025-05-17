package request

type UpdateInformationTenantRequest struct {
	Name    *string `json:"name"`
	LogoUrl *string `json:"logo_url"`
}

type UpdateStatusTenantRequest struct {
	IsActive *int `json:"is_active" validate:"required;numeric"`
}
