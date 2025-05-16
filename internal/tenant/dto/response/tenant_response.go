package response

type TenantResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	LogoUrl  string `json:"logo_url"`
	IsActive int    `json:"is_active"`
}

type TenantListItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
