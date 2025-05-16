package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/model"
)

func FromRoleModel(role *model.Role) *response.RoleResponse {
	return &response.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
	}
}

func FromRoleToListItem(role *model.Role) *response.RoleListItemResponse {
	return &response.RoleListItemResponse{
		ID:   role.ID,
		Name: role.Name,
		Slug: role.Slug,
	}
}

func FromRoleToList(roles []*model.Role) []*response.RoleListItemResponse {
	result := make([]*response.RoleListItemResponse, 0, len(roles))
	for _, role := range roles {
		result = append(result, FromRoleToListItem(role))
	}
	return result
}
