package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
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
	Create(req request.CreateUserRequest, auth middleware.AuthUser) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	logRepo  logRepository.LoginHistoryRepository
}

func NewUserUsecase(userRepo repository.UserRepository, roleRepo repository.RoleRepository, logRepo logRepository.LoginHistoryRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		logRepo:  logRepo,
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
	payload := map[string]interface{}{
		"id":        getUser.ID,
		"name":      getUser.Name,
		"email":     getUser.Email,
		"tenant":    getUser.Tenant,
		"role":      getUser.Role,
		"is_active": getUser.IsActive,
		"token":     token,
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
		AccessToken: token,
		Exp:         10 * 60,
		TokenType:   "Bearer",
		Username:    req.Email,
	}, nil

}

func (u *userUsecase) Create(req request.CreateUserRequest, auth middleware.AuthUser) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	role, errGetRole := u.roleRepo.GetRole("organizer")

	if errGetRole != nil {
		log.WithFields(log.Fields{
			"error": errGetRole,
		}).Error("failed to get Role")
		return errGetRole
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

	if errCreateUser := u.userRepo.Create(user); errCreateUser != nil {
		log.WithFields(log.Fields{
			"error": errCreateUser,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong")
	}
	return nil
}
