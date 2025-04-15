package router

import (
	"test/pck/auth/controllers"

	"test/pck/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, authController *controllers.AuthController) {
	app.Post("/sign-up", authController.SignUpUser)
	app.Post("/sign-in", authController.SignInUser)
	app.Get("/logout", middleware.DeserializeUser, authController.LogOutUser)

	app.Get("/users/me", middleware.DeserializeUser, authController.GetMeHandler)
	app.Get("/users/", middleware.DeserializeUser, middleware.AllowedRoles([]string{"admin", "moderator"}), authController.GetUsersHandler)
}
