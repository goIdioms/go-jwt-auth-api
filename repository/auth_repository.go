package repository

import (
	"context"
	"fmt"
	"test/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	SignUpUser(*models.User) (*models.User, error)
	SignInUser(*models.User) (*models.User, error)
}

type AuthRepositoryImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthRepository(collection *mongo.Collection, ctx context.Context) AuthRepository {
	return &AuthRepositoryImpl{collection, ctx}
}

func (r *AuthRepositoryImpl) SignUpUser(payload *models.User) (*models.User, error) {

	res, err := r.collection.InsertOne(r.ctx, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user models.User
	query := bson.M{"_id": res.InsertedID}
	err = r.collection.FindOne(r.ctx, query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) SignInUser(payload *models.User) (*models.User, error) {
	return nil, nil
}
