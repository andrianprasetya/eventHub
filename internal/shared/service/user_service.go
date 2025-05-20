package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	appServer "github.com/andrianprasetya/eventHub/server"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func MapToAuthUserPayload(user *modelUser.User, token string) *middleware.AuthUser {
	return &middleware.AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Tenant: middleware.TenantPayload{
			ID:       user.Tenant.ID,
			Name:     user.Tenant.Name,
			Email:    user.Tenant.Email,
			LogoUrl:  user.Tenant.LogoUrl,
			Domain:   user.Tenant.Domain,
			IsActive: user.Tenant.IsActive,
		},
		Role: middleware.RolePayload{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			Slug:        user.Role.Slug,
			Description: user.Role.Description,
		},
		IsActive: user.IsActive,
		Token:    token,
	}
}

func MapUserPayload(tenantID, roleID, name, email, password string) *modelUser.User {
	return &modelUser.User{
		ID:       utils.GenerateID(),
		TenantID: tenantID,
		RoleID:   roleID,
		Name:     name,
		Email:    email,
		Password: password,
		IsActive: 1,
	}
}

func BuildLoginResponse(authUser *middleware.AuthUser) *response.LoginResponse {
	return &response.LoginResponse{
		AccessToken:  authUser.Token,
		Exp:          10 * 60,
		TokenType:    "Bearer",
		Username:     authUser.Email,
		TenantDomain: authUser.Tenant.Domain,
	}
}

func HandleRedisToken(ctx context.Context, userID string) (*middleware.AuthUser, error) {
	tokenData, err := appServer.RedisClient.Get(ctx, "user:jwt:"+userID)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // not found is okay
		}
		return nil, err // redis error
	}

	var authUser middleware.AuthUser
	if err := json.Unmarshal([]byte(tokenData), &authUser); err != nil {
		return nil, err
	}
	return &authUser, nil
}

func SaveTokenToRedis(ctx context.Context, userID string, payload *middleware.AuthUser) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	key := "user:jwt:" + userID
	_, err = appServer.RedisClient.SetWithExpire(ctx, key, data, 10*time.Minute)
	return err
}

func CheckMaxUserCanCreated(countUserCreated int, tenantSetting *modelTenant.TenantSettingChannel) error {
	maxEvent, _ := strconv.Atoi(tenantSetting.TenantSetting.Value)

	if countUserCreated >= maxEvent {

		return validation.ValidationError{
			"event_limit": fmt.Sprintf("your subscription package not able create event more than %d", maxEvent),
		}
	}
	return nil
}
