package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/goIdioms/go-jwt-auth-api/pck/database"
	"github.com/goIdioms/go-jwt-auth-api/pck/models"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthRepository interface {
	SignUpUser(*models.User) (*models.User, error)
	SignInUser(*models.SignInInput) (*models.User, error)
}

type AuthRepositoryImpl struct {
	ctx context.Context
}

func NewAuthRepository(ctx context.Context) AuthRepository {
	return &AuthRepositoryImpl{ctx}
}

func (r *AuthRepositoryImpl) SignUpUser(payload *models.User) (*models.User, error) {
	res, err := database.UserCollection.InsertOne(r.ctx, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user models.User
	query := bson.M{"_id": res.InsertedID}
	err = database.UserCollection.FindOne(r.ctx, query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) SignInUser(payload *models.SignInInput) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": strings.ToLower(payload.Email)}
	err := database.UserCollection.FindOne(r.ctx, filter).Decode(&user)

	return &user, err
}
