package services

import (
	"fmt"
	"strings"
	"test/database"
	"test/models"
	"test/repository"
	"test/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	SignUpUser(payload *models.SignUpInput) (*models.User, error)
	SignInUser(payload *models.SignInInput) (string, error)
}

type AuthServiceImpl struct {
	userRepo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{userRepo: repo}
}

func (s *AuthServiceImpl) SignUpUser(payload *models.SignUpInput) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
		Photo:    payload.Photo,
	}
	result, err := s.userRepo.SignUpUser(&newUser)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return result, nil
}

func (s *AuthServiceImpl) SignInUser(payload *models.SignInInput) (string, error) {
	user, err := s.userRepo.SignInUser(payload)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return "", fmt.Errorf("invalid email or password")
		}
		return "", err
	}

	err = utils.CompareHashAndPassword(user.Password, payload.Password)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	config, _ := database.LoadConfig(".")
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID                             // идентификатор пользователя
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix() // время истечения срока действия токена
	claims["iat"] = now.Unix()                          // время создания токена
	claims["nbf"] = now.Unix()                          // время, до которого токен действителен
	claims["exp"] = now.Add(24 * time.Hour).Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	return tokenString, err
}
