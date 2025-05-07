package response

type UserLog struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	RoleID   string `json:"role_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive int    `json:"is_active"`
}
