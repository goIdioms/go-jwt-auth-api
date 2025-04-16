package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(ttl time.Duration, userID interface{}, privatKey string) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = userID
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := tokenByte.SignedString([]byte(privatKey))

	return token, err
}
