package usecase

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
)

type RoleUsecase interface {
	GetAll(page, pageSize int) ([]*response.RoleListItemResponse, int64, error)
	GetByID(id string) (*response.RoleResponse, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUsecase(roleRepo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{roleRepo: roleRepo}
}

func (r *roleUsecase) GetAll(page, pageSize int) ([]*response.RoleListItemResponse, int64, error) {
	roles, total, err := r.roleRepo.GetAll(page, pageSize)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get Role")
		return nil, 0, fmt.Errorf("something Went wrong %w", err)
	}

	return mapper.FromRoleToList(roles), total, err
}

func (r *roleUsecase) GetByID(id string) (*response.RoleResponse, error) {
	role, err := r.roleRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get Role ")
		return nil, fmt.Errorf("something Went wrong %w", err)
	}
	return mapper.FromRoleModel(role), err

}
