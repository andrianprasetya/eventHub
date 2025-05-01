package usecase

import (
	"context"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/user"
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/dto/response"
	appServer "github.com/andrianprasetya/eventHub/server"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {

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
	key := "jwt:" + token
	if errGenerateJwt != nil {
		log.WithFields(log.Fields{
			"error": errGenerateJwt,
		}).Error("failed to generate jwt")
		return nil, errGenerateJwt
	}
	_, errRedis := appServer.RedisClient.SetWithExpire(ctx, key, token, 10*time.Minute)
	if errRedis != nil {
		log.WithFields(log.Fields{
			"error": errRedis,
		}).Error("failed to save token in redis")
		return nil, errRedis
	}

	return &response.LoginResponse{
		AccessToken: token,
		Exp:         10 * 60,
		TokenType:   "Bearer",
		Username:    req.Email,
	}, nil

}
