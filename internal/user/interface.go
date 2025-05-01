package user

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
}

type RoleRepository interface {
	CreateUserRole(userRole *model.UserRole) error
	GetRole(name string) (*model.Role, error)
}

type UserUsecase interface {
	Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error)
}
