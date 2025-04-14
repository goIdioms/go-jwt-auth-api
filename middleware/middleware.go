package middleware

import (
	"context"
	"fmt"
	"strings"
	"test/database"
	"test/models"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Authorization token is missing",
		})
	}

	config, _ := database.LoadConfig(".")
	tokenByte, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalidate token"})
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

	if user.ID.Hex() != fmt.Sprint(claims["sub"]) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "The user belonging to this token no longer exists",
		})
	}

	c.Locals("user", models.FilteredUserResponse(&user))
	c.Next()
	return nil
}
