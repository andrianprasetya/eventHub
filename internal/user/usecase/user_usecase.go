package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	responseDTO "github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/user/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"github.com/andrianprasetya/eventHub/internal/user/repository"
	appServer "github.com/andrianprasetya/eventHub/server"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserUsecase interface {
	Login(ctx context.Context, req request.LoginRequest, ip string) (*response.LoginResponse, error)
	Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) error
	GetAll(page, pageSize int) ([]*response.UserListItemResponse, int64, error)
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
	getUser, errGetUser := u.userRepo.GetByEmail(req.Email)

	if errGetUser != nil {
		log.WithFields(log.Fields{
			"error": errGetUser,
		}).Error("failed to get Email")
		return nil, errGetUser
	}

	if errMatching := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(req.Password)); errMatching != nil {
		log.WithFields(log.Fields{
			"error": errMatching,
		}).Error("failed to matching password")
		return nil, errMatching
	}

	token, errGenerateJwt := utils.GenerateJWT(req.Email)

	payload := &middleware.AuthUser{
		ID:    getUser.ID,
		Name:  getUser.Name,
		Email: getUser.Email,
		Tenant: middleware.TenantPayload{
			ID:       getUser.Tenant.ID,
			Name:     getUser.Tenant.Name,
			Email:    getUser.Tenant.Email,
			LogoUrl:  getUser.Tenant.LogoUrl,
			Domain:   getUser.Tenant.Domain,
			IsActive: getUser.Tenant.IsActive,
		},
		Role: middleware.RolePayload{
			ID:          getUser.Role.ID,
			Name:        getUser.Role.Name,
			Slug:        getUser.Role.Slug,
			Description: getUser.Role.Description,
		},
		IsActive: getUser.IsActive,
		Token:    token,
	}
	data, _ := json.Marshal(payload)
	key := "user:jwt:" + token
	if errGenerateJwt != nil {
		log.WithFields(log.Fields{
			"error": errGenerateJwt,
		}).Error("failed to generate jwt")
		return nil, errGenerateJwt
	}
	_, errRedis := appServer.RedisClient.SetWithExpire(ctx, key, data, 10*time.Minute)
	if errRedis != nil {
		log.WithFields(log.Fields{
			"error": errRedis,
		}).Error("failed to save token in redis")
		return nil, errRedis
	}

	//login login time
	helper.LogLoginHistory(u.logRepo, getUser.ID, ip)

	return &response.LoginResponse{
		AccessToken:  token,
		Exp:          10 * 60,
		TokenType:    "Bearer",
		Username:     req.Email,
		TenantDomain: getUser.Tenant.Domain,
	}, nil

}

func (u *userUsecase) Create(req request.CreateUserRequest, auth *middleware.AuthUser, url string) error {
	var err error
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to bcrypt password")
		return err
	}

	role, err := u.roleRepo.GetByID("organizer")

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get Role")
		return err
	}

	user := &modelUser.User{
		ID:       utils.GenerateID(),
		TenantID: auth.Tenant.ID,
		RoleID:   role.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		IsActive: 1,
	}

	if err = u.userRepo.Create(user); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong %w", err)
	}

	userLog := responseDTO.UserLog{
		ID:       user.ID,
		TenantID: user.TenantID,
		RoleID:   user.RoleID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}

	userJSON, err := json.Marshal(userLog)
	if err != nil {
		return fmt.Errorf("error marshaling user")
	}

	if err == nil {
		helper.LogActivity(u.activityRepo, auth.ID, url, "Create User", string(userJSON), "user", user.ID)
	}

	return err
}

func (u *userUsecase) GetAll(page, pageSize int) ([]*response.UserListItemResponse, int64, error) {
	users, total, err := u.userRepo.GetAll(page, pageSize)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get users")
		return nil, 0, fmt.Errorf("something Went wrong %w", err)
	}
	return mapper.FromUserToList(users), total, err
}

func (u *userUsecase) GetByID(id string) (*response.UserResponse, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get users")
		return nil, fmt.Errorf("something Went wrong %w", err)
	}
	return mapper.FromUserModel(user), err
}

func (u *userUsecase) Update(req request.UpdateUserRequest, id string) error {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get user")
		return fmt.Errorf("something Went wrong %w", err)
	}

	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := u.userRepo.Update(user); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong %w", err)
	}

	return nil

}

func (u *userUsecase) Delete(id string) error {
	if err := u.userRepo.Delete(id); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to delete user")
		return fmt.Errorf("something Went wrong %w", err)
	}
	return nil
}
