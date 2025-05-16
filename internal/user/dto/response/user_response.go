package response

import (
	tenantResponse "github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
)

type UserResponse struct {
	ID       string                        `json:"id"`
	Tenant   tenantResponse.TenantResponse `json:"tenant"`
	Role     RoleResponse                  `json:"role"`
	Name     string                        `json:"name"`
	Email    string                        `json:"email"`
	IsActive int                           `json:"is_active"`
}

type UserListItemResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RoleName string `json:"role_name"`
	IsActive int    `json:"is_active"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	Username     string `json:"username"`
	TokenType    string `json:"token_type"`
	Exp          int64  `json:"exp"`
	TenantDomain string `json:"tenant_domain"`
}
