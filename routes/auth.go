package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupAuthRoutes(app *fiber.App) {
	authService := service.NewAuthService()
	authController := controller.NewAuthController(authService)

	auth := app.Group("/api/auth")

	auth.Post("/Register", authController.Register)
	auth.Post("/Login", authController.Login)

	auth.Post("/Logout", middleware.AuthMiddleware(), authController.Logout)
}

func SetupProtectedRoutes(app *fiber.App) {
	api := app.Group("/api", middleware.AuthMiddleware())

	api.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		userEmail := c.Locals("user_email").(string)

		return c.JSON(fiber.Map{
			"message":    "This is a protected route",
			"user_id":    userID,
			"user_email": userEmail,
		})
	})
}
