package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupProfileRoutes(app *fiber.App) {
	profileService := service.NewProfileService()
	profileController := controller.NewProfileController(profileService)

	profile := app.Group("/api", middleware.AuthMiddleware()) //
	profile.Get("/profile", profileController.GetUserProfile)
	profile.Put("/update-profile", profileController.UpdateProfile)
}
