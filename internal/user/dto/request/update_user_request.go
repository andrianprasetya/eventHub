package request

type UpdateUserRequest struct {
	RoleID   *string `json:"role_id"`
	IsActive *int    `json:"is_active"`
}
