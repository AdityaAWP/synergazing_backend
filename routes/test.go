package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
)

func SetupTestRoutes(app *fiber.App) {
	testController := controller.NewTestController()

	test := app.Group("/test")

	// Create test users
	test.Post("/users", testController.CreateTestUsers)

	// Create test chat
	test.Post("/chat", testController.CreateTestChat)

	// Get test token for a user
	test.Get("/token/:user_id", testController.GetTestToken)

	// List all test data
	test.Get("/data", testController.ListTestData)
}
