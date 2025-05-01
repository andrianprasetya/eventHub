package response

type TenantResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TenantListItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
