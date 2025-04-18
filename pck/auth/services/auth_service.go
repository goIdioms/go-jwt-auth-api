package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/goIdioms/go-jwt-auth-api/pck/auth/repository"
	"github.com/goIdioms/go-jwt-auth-api/pck/cache"
	"github.com/goIdioms/go-jwt-auth-api/pck/database"
	"github.com/goIdioms/go-jwt-auth-api/pck/models"
	"github.com/goIdioms/go-jwt-auth-api/pck/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	SignUpUser(payload *models.SignUpInput) (*models.User, error)
	SignInUser(ctx *fiber.Ctx, payload *models.SignInInput) (*models.Tokens, error)
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

func (s *AuthServiceImpl) SignInUser(c *fiber.Ctx, payload *models.SignInInput) (*models.Tokens, error) {
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

	accessToken, err := utils.GenerateToken(ttlAccess, user.ID, scrAccess)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %v", err)
	}

	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		existingRefreshToken, err := s.cache.GetRefreshToken(c.Context(), sessionID)
		if err != nil {
			return nil, fmt.Errorf("error getting refresh token: %v", err)
		}
		if existingRefreshToken != nil {
			return &models.Tokens{
				SessionID:    sessionID,
				AccessToken:  accessToken,
				RefreshToken: existingRefreshToken.RefreshToken,
			}, nil
		}
	}

	newRefreshToken, err := utils.GenerateToken(ttlRefresh, user.ID, scrRefresh)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %v", err)
	}

	sessionID = uuid.New().String()
	value := cache.CacheValue{
		UserID:       user.ID.Hex(),
		RefreshToken: newRefreshToken,
	}
	err = s.cache.SaveRefreshToken(s.ctx, sessionID, value, time.Duration(ttlRefresh))
	fmt.Println("save to cache")
	if err != nil {
		return nil, fmt.Errorf("error saving refresh token: %v", err)
	}

	tokens := &models.Tokens{
		SessionID:    sessionID,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	return tokens, nil
}
