package request

type UpdateTenantRequest struct {
	Name    string `json:"name" validate:"required"`
	LogoUrl string `json:"logo_url" validate:"required"`
}
