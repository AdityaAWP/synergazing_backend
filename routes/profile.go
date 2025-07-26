package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/handler"
	"synergazing.com/synergazing/middleware"
)

func SetupProfileRoutes(app *fiber.App)  {
	profile := app.Group("/api/profile", middleware.AuthMiddleware())

	profile.Get("/", handler.GetUserProfile)
}