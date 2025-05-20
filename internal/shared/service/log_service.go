package service

import (
	responseDTO "github.com/andrianprasetya/eventHub/internal/shared/response"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
)

func MapUserLog(user *modelUser.User) *responseDTO.UserLog {
	return &responseDTO.UserLog{
		ID:       user.ID,
		TenantID: user.TenantID,
		RoleID:   user.RoleID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}
