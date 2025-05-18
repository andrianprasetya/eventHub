package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var SecretKey = os.Getenv("JWT_SECRET")

func GenerateJWT(userID, username string) (string, error) {
	if SecretKey == "" {
		log.WithFields(log.Fields{
			"errors": "cannot get value secret key or not declared on env",
		}).Error("secret key is null for jwt")
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func DecodeJWT(tokenString string) (map[string]interface{}, error) {
	// Parse token dan verifikasi dengan secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Cek algoritma signing-nya
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Ambil claims jika token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Ubah ke map[string]interface{} agar fleksibel
		result := make(map[string]interface{})
		for key, val := range claims {
			result[key] = val
		}
		return result, nil
	}

	return nil, errors.New("could not parse claims")
}
