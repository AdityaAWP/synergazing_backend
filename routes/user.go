package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
)

func SetupUserRoutes(app *fiber.App) {
	users := app.Group("/api", middleware.AuthMiddleware())
	users.Get("/users", controller.ListAllUsers)
	users.Get("/users/ready", controller.GetReadyUsers)
}
