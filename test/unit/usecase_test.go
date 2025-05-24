package unit

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	appServer "github.com/andrianprasetya/eventHub/server"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
)

var u = modelUser.User{
	Email:    "andrianprasetya222@gmail.com",
	Password: "test1234",
	Name:     "Andrian Prasetya",
}

func TestHashedAndCheckPassword(t *testing.T) {
	t.Run("successfully hash and verify password", func(t *testing.T) {
		hashed, err := service.HashedPassword(u.Password)
		require.NoError(t, err)
		require.NotEmpty(t, hashed)

		err = service.CheckPassword(string(hashed), u.Password)
		require.NoError(t, err)
	})

	t.Run("fail verify wrong password", func(t *testing.T) {
		hashed, _ := service.HashedPassword(u.Password)

		err := service.CheckPassword(string(hashed), "wrongpassword")
		require.Error(t, err)
	})
}

func TestMapToAuthUserPayload(t *testing.T) {
	user := &modelUser.User{
		ID:    "uid123",
		Name:  "Andrian",
		Email: "test@email.com",
		Tenant: modelTenant.Tenant{
			ID:       "tid",
			Name:     "TenantX",
			Email:    "tenant@email.com",
			LogoUrl:  "logo.png",
			Domain:   "tenant.com",
			IsActive: 1,
		},
		Role: modelUser.Role{
			ID:          "rid",
			Name:        "Admin",
			Slug:        "admin",
			Description: "Administrator",
		},
		IsActive: 1,
	}

	token := "jwt.token.here"
	result := service.MapToAuthUserPayload(user, token)

	require.Equal(t, token, result.Token)
	require.Equal(t, user.ID, result.ID)
	require.Equal(t, user.Tenant.Name, result.Tenant.Name)
}

func TestMapUserPayload(t *testing.T) {
	u := service.MapUserPayload("tid", "rid", "Andrian", "test@email.com", "pass123")

	require.Equal(t, "tid", u.TenantID)
	require.Equal(t, "rid", u.RoleID)
	require.Equal(t, "Andrian", u.Name)
	require.Equal(t, "test@email.com", u.Email)
	require.Equal(t, "pass123", u.Password)
	require.Equal(t, 1, u.IsActive)
}

func TestBuildLoginResponse(t *testing.T) {
	authUser := &middleware.AuthUser{
		Email: "user@email.com",
		Token: "access.token",
		Tenant: middleware.TenantPayload{
			Domain: "tenant.domain.com",
		},
	}

	resp := service.BuildLoginResponse(authUser)
	require.Equal(t, "access.token", resp.AccessToken)
	require.Equal(t, "user@email.com", resp.Username)
	require.Equal(t, "tenant.domain.com", resp.TenantDomain)
	require.Equal(t, "Bearer", resp.TokenType)
	require.Equal(t, int64(600), resp.Exp)
}

func TestHandleAndSaveTokenToRedis(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	appServer.RedisClient = redisser.NewRedisClient(rdb)

	ctx := context.Background()
	authUser := &middleware.AuthUser{ID: "uid", Email: "user@email.com"}

	// Test Save
	err := service.SaveTokenToRedis(ctx, "uid", authUser)
	require.NoError(t, err)

	// Test Load
	result, err := service.HandleRedisToken(ctx, "uid")
	require.NoError(t, err)
	require.Equal(t, "user@email.com", result.Email)

	// Test Not Found
	res, err := service.HandleRedisToken(ctx, "unknown")
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestCheckMaxUserCanCreated(t *testing.T) {
	t.Run("allowed to create user", func(t *testing.T) {
		setting := &modelTenant.TenantSetting{Value: "5"}
		err := service.CheckMaxUserCanCreated(4, setting)
		require.NoError(t, err)
	})

	t.Run("exceeds user quota", func(t *testing.T) {
		setting := &modelTenant.TenantSetting{Value: "3"}
		err := service.CheckMaxUserCanCreated(5, setting)
		require.Error(t, err)
		require.Contains(t, err.Error(), "created user")
	})
}
