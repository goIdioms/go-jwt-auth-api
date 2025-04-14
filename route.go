package main

import (
	"test/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, authController *controllers.AuthController) {
	app.Post("/sign-up", authController.SignUpUser)
	app.Post("/sign-in", authController.SignInUser)
	app.Get("/logout", DeserializeUser, authController.LogOutUser)

	app.Get("/users/me", DeserializeUser, authController.GetMeHandler)
	app.Get("/users/", DeserializeUser, AllowedRoles([]string{"admin", "moderator"}), authController.GetUsersHandler)
}
