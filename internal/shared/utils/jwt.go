package utils

import (
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var SecretKey = os.Getenv("JWT_SECRET")

func GenerateJWT(username string) (string, error) {
	if SecretKey == "" {
		log.WithFields(log.Fields{
			"error": "cannot get value secret key or not declared on env",
		}).Error("secret key is null for jwt")
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
