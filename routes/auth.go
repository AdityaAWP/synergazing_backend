package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/service"
)

func SetupAuthRoutes(app *fiber.App) {
	db := config.GetDB()
	otpService := service.NewOTPService()
	authService := service.NewAuthService(otpService)
	socialAuthService := service.NewSocialAuthService(db)

	authController := controller.NewAuthController(authService, otpService)
	socialController := controller.NewSocialController(socialAuthService, authService)

	auth := app.Group("/api/auth")

	auth.Post("/register/initiate", authController.InitiateRegistration)
	auth.Post("/register/complete", authController.CompleteRegistration)
	auth.Post("/otp/resend", authController.ResendOTP)
	auth.Post("/otp/verify", authController.VerifyOTP)

	auth.Get("/register/info", authController.GetRegistrationInfo)
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)

	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)

	// Email verification endpoints
	auth.Post("/verify-email", authController.VerifyEmail)
	auth.Post("/request-email-verification", authController.RequestEmailVerification)
	auth.Get("/verification-status", authController.GetUserVerificationStatus)

	google := auth.Group("/google")
	google.Get("/login", socialController.GoogleLogin)
	google.Get("/callback", socialController.GoogleCallback)
}
