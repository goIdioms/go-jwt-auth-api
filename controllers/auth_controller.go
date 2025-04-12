package controllers

import (
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
	var payload *models.SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"stattus": "fail",
			"message": err.Error(),
		})
	}

	erroros := models.ValidateStruct(payload)
	if erroros != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"stattus": "fail",
			"message": "Password do not match",
		})
	}

	user, err := ac.AuthService.SignUpUser(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"stattus": "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stattus": "success",
		"user":    user,
	})
}
