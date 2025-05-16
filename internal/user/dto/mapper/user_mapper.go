package mapper

import (
	tenantResponse "github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	userResponse "github.com/andrianprasetya/eventHub/internal/user/dto/response"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
)

func FromUserModel(user *modelUser.User) *userResponse.UserResponse {
	return &userResponse.UserResponse{
		ID: user.ID,
		Tenant: tenantResponse.TenantResponse{
			ID:       user.Tenant.ID,
			Name:     user.Tenant.Name,
			Email:    user.Tenant.Email,
			LogoUrl:  user.Tenant.LogoUrl,
			IsActive: user.Tenant.IsActive,
		},
		Role: userResponse.RoleResponse{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			Slug:        user.Role.Slug,
			Description: user.Role.Description,
		},
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}

func FromUserToListItem(user *modelUser.User, roleName string) *userResponse.UserListItemResponse {
	return &userResponse.UserListItemResponse{
		ID:       user.ID,
		Name:     user.Name,
		RoleName: roleName,
		IsActive: user.IsActive,
	}
}

func FromUserToList(users []*modelUser.User) []*userResponse.UserListItemResponse {
	result := make([]*userResponse.UserListItemResponse, 0, len(users))
	for _, user := range users {
		result = append(result, FromUserToListItem(user, user.Role.Name))
	}
	return result
}
