package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupChatRoutes(app *fiber.App) {
	chatService := service.NewChatService()
	chatController := controller.NewChatController(chatService)

	// WebSocket route (no auth middleware for WebSocket upgrade)
	app.Use("/ws/chat", chatController.WebSocketUpgrade)
	app.Get("/ws/chat", websocket.New(chatController.HandleWebSocket))

	// REST API routes (protected)
	api := app.Group("/api/chat", middleware.AuthMiddleware())

	// Get or create chat with another user
	api.Get("/with/:user_id", chatController.GetOrCreateChat)

	// Get all chats for current user
	api.Get("/", chatController.GetUserChats)

	// Get messages for a specific chat
	api.Get("/:chat_id/messages", chatController.GetChatMessages)

	// Mark chat messages as read
	api.Put("/:chat_id/read", chatController.MarkChatAsRead)
}
