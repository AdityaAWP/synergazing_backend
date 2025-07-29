package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
)

func SetupProfileRoutes(app *fiber.App) {
	profile := app.Group("/api", middleware.AuthMiddleware())
	profile.Get("/profile", controller.GetUserProfile)
	profile.Put("/update-profile", controller.UpdateProfile)
}
