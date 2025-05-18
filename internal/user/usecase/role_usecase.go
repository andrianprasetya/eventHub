package usecase

import (
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
)

type RoleUsecase interface {
	GetAll(query request.RolePaginateParams) ([]*response.RoleListItemResponse, int64, error)
	GetByID(id string) (*response.RoleResponse, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUsecase(roleRepo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{roleRepo: roleRepo}
}

func (r *roleUsecase) GetAll(query request.RolePaginateParams) ([]*response.RoleListItemResponse, int64, error) {
	roles, total, err := r.roleRepo.GetAll(query)

	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get Role")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromRoleToList(roles), total, err
}

func (r *roleUsecase) GetByID(id string) (*response.RoleResponse, error) {
	role, err := r.roleRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get Role")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromRoleModel(role), err

}
