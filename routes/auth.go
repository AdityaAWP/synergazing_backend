package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/service"
)

func SetupAuthRoutes(app *fiber.App) {
	db := config.GetDB()
	authService := service.NewAuthService()
	socialAuthService := service.NewSocialAuthService(db)

	authController := controller.NewAuthController(authService)
	socialController := controller.NewSocialController(socialAuthService, authService)

	auth := app.Group("/api/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)

	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)

	google := auth.Group("/google")
	google.Get("/login", socialController.GoogleLogin)
	google.Get("/callback", socialController.GoogleCallback)
}
