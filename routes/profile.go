package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/handler"
	"synergazing.com/synergazing/middleware"
)

func SetupProfileRoutes(app *fiber.App) {
	profile := app.Group("/api/profile", middleware.AuthMiddleware())
	fmt.Println("asdasdasd")
	profile.Get("/", handler.GetUserProfile)
	profile.Put("/update-profile", handler.UpdateProfile)
}
