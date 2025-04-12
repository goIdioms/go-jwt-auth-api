package services

import "test/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.UserResponse, error)
	SignInUser(*models.SignInInput) (*models.UserResponse, error)
}
