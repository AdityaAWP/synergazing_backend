package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupNotificationRoutes(app *fiber.App) {
	db := config.GetDB()
	notificationService := service.NewNotificationService(db)
	notificationController := controller.NewNotificationController(notificationService)

	notifications := app.Group("/api/notifications", middleware.AuthMiddleware())

	notifications.Get("/", notificationController.GetNotifications)

	notifications.Get("/unread", notificationController.GetUnreadNotifications)

	notifications.Get("/count", notificationController.GetUnreadCount)

	notifications.Put("/:id/read", notificationController.MarkAsRead)

	notifications.Put("/read-all", notificationController.MarkAllAsRead)

	notifications.Delete("/:id", notificationController.DeleteNotification)

	notifications.Post("/test-deadlines", notificationController.TestDeadlineNotifications)
}
