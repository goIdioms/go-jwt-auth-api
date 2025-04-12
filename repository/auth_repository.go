package repository

import (
	"test/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	SignUpUser(*models.SignUpInput) (*models.UserResponse, error)
	SignInUser(*models.SignInInput) (*models.UserResponse, error)
}

type AuthRepositoryImpl struct {
	collection *mongo.Collection
}

func NewAuthRepository(collection *mongo.Collection) AuthRepository {
	return &AuthRepositoryImpl{collection}
}

func (r *AuthRepositoryImpl) SignUpUser(input *models.SignUpInput) (*models.UserResponse, error) {
	return nil, nil
}

func (r *AuthRepositoryImpl) SignInUser(input *models.SignInInput) (*models.UserResponse, error) {
	return nil, nil
}
