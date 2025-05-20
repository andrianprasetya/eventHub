package usecase

import (
	"context"
	"encoding/json"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
)

type UserUsecase interface {
	Login(ctx context.Context, req request.LoginRequest, ip string) (*response.LoginResponse, error)
	Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) error
	GetAll(query request.UserPaginateParams, tenantID *string) ([]*response.UserListItemResponse, int64, error)
	GetByID(id string) (*response.UserResponse, error)
	Update(req request.UpdateUserRequest, id string) error
	Delete(id string) error
}

type userUsecase struct {
	txManager    repositoryShared.TxManager
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	logRepo      logRepository.LoginHistoryRepository
	activityRepo logRepository.LogActivityRepository
}

func NewUserUsecase(
	txManager repositoryShared.TxManager,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	logRepo logRepository.LoginHistoryRepository,
	activityRepo logRepository.LogActivityRepository) UserUsecase {
	return &userUsecase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		logRepo:      logRepo,
		activityRepo: activityRepo,
		txManager:    txManager,
	}
}

func (u *userUsecase) Login(ctx context.Context, req request.LoginRequest, ip string) (*response.LoginResponse, error) {
	getUser, err := u.userRepo.GetByEmail(req.Email)

	//data tidak ada
	if err != nil || getUser == nil {
		log.WithFields(log.Fields{
			"errors": err,
			"email":  req.Email,
		}).Error("failed to get email")
		return nil, appErrors.ErrInvalidCredentials
	}

	// Cek password
	if err := service.CheckPassword(getUser.Password, req.Password); err != nil {
		log.WithError(err).Error("invalid password")
		return nil, appErrors.ErrInvalidCredentials
	}

	// Cek token redis
	if cachedUser, err := service.HandleRedisToken(ctx, getUser.ID); err != nil {
		log.WithError(err).Error("failed to get token from redis")
		return nil, appErrors.ErrInternalServer
	} else if cachedUser != nil {
		return service.BuildLoginResponse(cachedUser), nil
	}

	// Generate JWT dan mapping payload
	token, err := utils.GenerateJWT(getUser.ID, req.Email)
	if err != nil {
		log.WithError(err).Error("failed to generate jwt")
		return nil, appErrors.ErrInternalServer
	}

	authPayload := service.MapToAuthUserPayload(getUser, token)

	// Simpan ke Redis
	if err := service.SaveTokenToRedis(ctx, getUser.ID, authPayload); err != nil {
		log.WithError(err).Error("failed to save token to redis")
		return nil, appErrors.ErrInternalServer
	}

	// Log history login
	helper.LogLoginHistory(u.logRepo, getUser.ID, ip)
	return service.BuildLoginResponse(authPayload), nil
}

func (u *userUsecase) Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) error {
	var err error

	hashedPassword, err := service.HashedPassword(req.Password)

	if err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"password": req.Password,
		}).Error("failed to bcrypt password")
		return appErrors.ErrInternalServer
	}

	role, err := u.roleRepo.GetByID(req.RoleID)

	if err != nil {
		log.WithFields(log.Fields{
			"errors":  err,
			"role_id": req.RoleID,
		}).Error("failed to get role")
		return appErrors.ErrInternalServer
	}

	user := service.MapUserPayload(auth.Tenant.ID, role.ID, req.Name, req.Email, string(hashedPassword))

	if err = u.userRepo.Create(user); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"user":   user,
		}).Error("failed to create user")
		return appErrors.ErrInternalServer
	}

	userLog := service.MapUserLog(user)

	userJSON, err := json.Marshal(userLog)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"user_log": userLog,
		}).Error("failed to marshal json user log")
		return appErrors.ErrInternalServer
	}

	helper.LogActivity(u.activityRepo, auth.Tenant.ID, auth.ID, url, "Create User", string(userJSON), "user", user.ID)

	return nil
}

func (u *userUsecase) GetAll(query request.UserPaginateParams, tenantID *string) ([]*response.UserListItemResponse, int64, error) {
	users, total, err := u.userRepo.GetAll(query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    tenantID,
		}).Error("failed to get users")
		return nil, 0, appErrors.ErrInternalServer
	}
	return mapper.FromUserToList(users), total, nil
}

func (u *userUsecase) GetByID(id string) (*response.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to get user")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromUserModel(user), nil
}

func (u *userUsecase) Update(req request.UpdateUserRequest, id string) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to get user")
		return appErrors.ErrInternalServer
	}

	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := u.userRepo.Update(user); err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"req_user": req,
		}).Error("failed to update user")
		return appErrors.ErrInternalServer
	}

	return nil

}

func (u *userUsecase) Delete(id string) error {
	if err := u.userRepo.Delete(id); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to delete user")
		return appErrors.ErrInternalServer
	}
	return nil
}
