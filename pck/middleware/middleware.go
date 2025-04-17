package middleware

import (
	"context"
	"fmt"
	"test/pck/cache"
	"test/pck/database"
	"test/pck/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeserializeUser(c *fiber.Ctx) error {
	// 1. Получаем refresh_token из куки
	sessionID := c.Cookies("session_id")
	redisCache := cache.NewRedisCache("localhost:6379")
	cachedValue, err := redisCache.GetRefreshToken(c.Context(), sessionID)
	if err != nil || cachedValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid or expired session",
		})
	}

	if sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Refresh token is missing",
		})
	}

	config, err := database.LoadConfig(".")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to load config",
		})
	}

	token, err := jwt.Parse(cachedValue.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.RefreshJwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid or expired refresh token",
		})
	}

	// 4. Получаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token claims",
		})
	}

	// 5. Получаем user из БД по ID
	userID := fmt.Sprint(claims["sub"])
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token subject",
		})
	}

	var user models.User
	err = database.UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
		})
	}

	// 6. Устанавливаем пользователя в контекст
	c.Locals("user", &user)
	return c.Next()
}

func AllowedRoles(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. User not authenticated",
			})
		}
		allowed := false
		for _, role := range allowedRoles {
			if role == user.Role {
				allowed = true
				break
			}
		}
		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. You are not allowed to perform this action",
			})
		}

		return c.Next()
	}

}
