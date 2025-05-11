package request

type UpdateTenantRequest struct {
	Name    *string `json:"name"`
	LogoUrl *string `json:"logo_url"`
}
