package services

import (
	"fmt"
	"strings"
	"test/internal/auth/repository"
	"test/internal/models"
	jwt "test/pkg/jwt"
	utils "test/pkg/security"

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

	tokenString, err := jwt.GenerateToken(user)

	return tokenString, err
}
