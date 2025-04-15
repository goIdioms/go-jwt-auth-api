package utils

import (
	"test/internal/database"
	"test/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *models.User) (string, error) {
	config, _ := database.LoadConfig(".")
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := tokenByte.SignedString([]byte(config.JwtSecret))
	return token, err
}
