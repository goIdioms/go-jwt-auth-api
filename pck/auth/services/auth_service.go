package services

import (
	"context"
	"fmt"
	"strings"
	"test/pck/auth/repository"
	"test/pck/cache"
	"test/pck/database"
	"test/pck/models"
	"test/pck/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	SignUpUser(payload *models.SignUpInput) (*models.User, error)
	SignInUser(payload *models.SignInInput) (*models.Tokens, error)
}

type AuthServiceImpl struct {
	userRepo repository.AuthRepository
	ctx      context.Context
	cache    *cache.RedisCache
}

func NewAuthService(ctx context.Context, repo repository.AuthRepository, cache *cache.RedisCache) AuthService {
	return &AuthServiceImpl{ctx: ctx, userRepo: repo, cache: cache}
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

func (s *AuthServiceImpl) SignInUser(payload *models.SignInInput) (*models.Tokens, error) {
	user, err := s.userRepo.SignInUser(payload)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	err = utils.CompareHashAndPassword(user.Password, payload.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	config, _ := database.LoadConfig(".")
	ttlAccess := config.AccessJwtExpiresIn
	ttlRefresh := config.RefreshJwtExpiresIn
	scrAccess := config.AccessJwtSecret
	scrRefresh := config.RefreshJwtSecret

	access_token, err := utils.GenerateToken(ttlAccess, user.ID, scrAccess)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %v", err)
	}

	refresh_token, err := utils.GenerateToken(ttlRefresh, user.ID, scrRefresh)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %v", err)
	}
	err = s.cache.Set(s.ctx, "01", string(refresh_token), 0)
	if err != nil {
		panic("failed to set refresh token in cache")
	}

	tokens := &models.Tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return tokens, nil
}
