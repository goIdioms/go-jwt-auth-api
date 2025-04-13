package controllers

import (
	"fmt"
	"test/models"
	"test/services"

	"github.com/gofiber/fiber/v2"
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
	// if err != nil {
	// 	if strings.Contains(err.Error(), "duplicate key") {
	// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
	// 			"status":  "fail",
	// 			"message": "User with that email already exists",
	// 		})
	// 	}
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": err.Error(),
	// 	})
	// }

	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "User creation failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user":   models.FilteredUserResponse(user),
	})
}
