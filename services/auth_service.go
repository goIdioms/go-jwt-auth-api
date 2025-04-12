package services

import (
	"context"
	"test/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.UserResponse, error)
	SignInUser(*models.SignInInput) (*models.UserResponse, error)
}

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) *AuthServiceImpl {
	return &AuthServiceImpl{
		collection: collection,
		ctx:        ctx,
	}
}
