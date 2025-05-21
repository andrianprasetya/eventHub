package usecase

import (
	"context"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

type RoleUsecase interface {
	GetAll(query request.RolePaginateParams, userAuth middleware.AuthUser) ([]*response.RoleListItemResponse, int64, error)
	GetByID(id string) (*response.RoleResponse, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUsecase(roleRepo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{roleRepo: roleRepo}
}

func (r *roleUsecase) GetAll(query request.RolePaginateParams, userAuth middleware.AuthUser) ([]*response.RoleListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	roles, total, err := r.roleRepo.GetAll(ctx, query, userAuth.Role.Slug)

	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"role":         userAuth.Role.Slug,
			"query_params": query,
		}).Error("failed to get role")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromRoleToList(roles), total, nil
}

func (r *roleUsecase) GetByID(id string) (*response.RoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to get Role")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromRoleModel(role), nil

}
