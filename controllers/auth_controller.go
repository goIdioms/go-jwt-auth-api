package controllers

import (
	"context"
	"fmt"
	"strconv"
	"test/database"
	"test/models"
	"test/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{AuthService: service}
}

func (ac *AuthController) SignUpUser(c *fiber.Ctx) error {
	payload := new(models.SignUpInput)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"stattus": "fail",
			"message": err.Error(),
		})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": errors,
		})
	}

	user, err := ac.AuthService.SignUpUser(payload)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user":   models.FilteredUserResponse(user),
	})
}

func (ac *AuthController) SignInUser(c *fiber.Ctx) error {
	payload := new(models.SignInInput)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"stattus": "fail",
			"message": err.Error(),
		})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": errors,
		})
	}

	token, err := ac.AuthService.SignInUser(payload)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("generating JWT Token failed: %v", err),
		})
	}

	config, _ := database.LoadConfig(".")
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   config.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}

func (ac *AuthController) LogOutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (ac *AuthController) GetMeHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user":   models.FilteredUserResponse(user),
	})
}

func (ac *AuthController) GetUsersHandler(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	cursor, err := database.UserCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch users",
		})
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse users",
		})
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.FilteredUserResponse(&user))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"results": len(userResponses),
		"users":   userResponses,
	})
}
