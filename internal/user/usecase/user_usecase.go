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
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	repositoryTenant "github.com/andrianprasetya/eventHub/internal/tenant/repository"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type UserUsecase interface {
	Login(req request.LoginRequest, ip string) (*response.LoginResponse, error)
	Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) error
	GetAll(query request.UserPaginateParams, tenantID *string) ([]*response.UserListItemResponse, int64, error)
	GetByID(id string) (*response.UserResponse, error)
	Update(req request.UpdateUserRequest, id string) error
	Delete(id string) error
}

type userUsecase struct {
	txManager         repositoryShared.TxManager
	userRepo          repository.UserRepository
	roleRepo          repository.RoleRepository
	tenantSettingRepo repositoryTenant.TenantSettingRepository
	logRepo           logRepository.LoginHistoryRepository
	activityRepo      logRepository.LogActivityRepository
}

func NewUserUsecase(
	txManager repositoryShared.TxManager,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantSettingRepo repositoryTenant.TenantSettingRepository,
	logRepo logRepository.LoginHistoryRepository,
	activityRepo logRepository.LogActivityRepository) UserUsecase {
	return &userUsecase{
		userRepo:          userRepo,
		roleRepo:          roleRepo,
		logRepo:           logRepo,
		tenantSettingRepo: tenantSettingRepo,
		activityRepo:      activityRepo,
		txManager:         txManager,
	}
}

func (u *userUsecase) Login(req request.LoginRequest, ip string) (*response.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	getUser, err := u.userRepo.GetByEmail(ctx, req.Email)

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

func (u *userUsecase) Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var (
		tenantSetting    *modelTenant.TenantSetting
		countUserCreated int
		unlimitedUser    string
		role             *modelUser.Role
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		getUnlimitedUser, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "unlimited_users")
		if err != nil {
			return err
		}
		unlimitedUser = getUnlimitedUser.Value
		return nil
	})

	g.Go(func() error {
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(ctx, auth.Tenant.ID, "max_users")
		if err != nil {
			return err
		}
		tenantSetting = getTenantSetting
		return nil
	})

	g.Go(func() error {
		getUserHasCreated, err := u.userRepo.CountCreatedUser(ctx, auth.Tenant.ID)
		if err != nil {
			return err
		}
		countUserCreated = getUserHasCreated
		return nil
	})

	g.Go(func() error {
		getRole, err := u.roleRepo.GetByID(ctx, req.RoleID)
		if err != nil {
			return err
		}
		role = getRole
		return nil
	})

	if err = g.Wait(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to fetch plan or role")
		return appErrors.WrapExpose(err, "failed to fetch plan or role", http.StatusInternalServerError)
	}

	if unlimitedUser == "false" {
		if err = service.CheckMaxUserCanCreated(countUserCreated, tenantSetting); err != nil {
			return appErrors.WrapExpose(err, "Created user quota Has been limit", http.StatusUnprocessableEntity)
		}
	}

	hashedPassword, err := service.HashedPassword(req.Password)

	if err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"password": req.Password,
		}).Error("failed to bcrypt password")
		return appErrors.ErrInternalServer
	}

	user := service.MapUserPayload(auth.Tenant.ID, role.ID, req.Name, req.Email, string(hashedPassword))

	if err = u.userRepo.Create(ctx, user); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"user":   user,
		}).Error("failed to create user")
		return appErrors.ErrInternalServer
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"user_log": user,
		}).Error("failed to marshal json user log")
		return appErrors.ErrInternalServer
	}

	helper.LogActivity(u.activityRepo, ctx, auth.Tenant.ID, auth.ID, url, "Create User", string(userJSON), "user", user.ID)

	return nil
}

func (u *userUsecase) GetAll(query request.UserPaginateParams, tenantID *string) ([]*response.UserListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	users, total, err := u.userRepo.GetAll(ctx, query, tenantID)
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to get user")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromUserModel(user), nil
}

func (u *userUsecase) Update(req request.UpdateUserRequest, id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	user, err := u.userRepo.GetByID(ctx, id)
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

	if err := u.userRepo.Update(ctx, user); err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"req_user": req,
		}).Error("failed to update user")
		return appErrors.ErrInternalServer
	}

	return nil

}

func (u *userUsecase) Delete(id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = u.userRepo.Delete(ctx, id); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to delete user")
		return appErrors.ErrInternalServer
	}
	return nil
}
