package services

import (
	"test/models"
	"test/repository"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.UserResponse, error)
	SignInUser(*models.SignInInput) (*models.UserResponse, error)
}

type AuthServiceImpl struct {
	userRepo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{userRepo: repo}
}

func (s *AuthServiceImpl) SignUpUser(input *models.SignUpInput) (*models.UserResponse, error) {
	return s.userRepo.SignUpUser(input)
}

func (s *AuthServiceImpl) SignInUser(input *models.SignInInput) (*models.UserResponse, error) {
	return s.userRepo.SignInUser(input)
}
