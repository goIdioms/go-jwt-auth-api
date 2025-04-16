package middleware

import (
	"context"
	"fmt"
	"test/pck/database"
	"test/pck/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeserializeUser(c *fiber.Ctx) error {
	tokenString := c.Cookies("refresh_token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Refresh token is missing",
		})
	}

	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token format",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid refresh token claims",
		})
	}

	var user models.User
	objectID, err := primitive.ObjectIDFromHex(fmt.Sprint(claims["sub"]))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token subject",
		})
	}

	filter := bson.M{"_id": objectID}
	err = database.UserCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "The user belonging to this token no longer exists",
		})
	}

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
