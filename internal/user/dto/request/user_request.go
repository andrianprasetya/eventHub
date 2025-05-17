package request

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	RoleID   *string `json:"role_id"`
	IsActive *int    `json:"is_active"`
}

type UserPaginateParams struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `query:"name"`
	Email    *string `query:"email"`
	RoleID   *string `query:"role_id"`
	TenantID *string `query:"tenant_id"`
	IsActive *int    `query:"is_active"`
}
