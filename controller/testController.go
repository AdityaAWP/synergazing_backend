package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
	"synergazing.com/synergazing/service"
)

type TestController struct {
	AuthService *service.AuthService
	ChatService *service.ChatService
}

func NewTestController() *TestController {
	return &TestController{
		AuthService: service.NewAuthService(),
		ChatService: service.NewChatService(),
	}
}

// CreateTestUsers creates test users for WebSocket testing
func (ctrl *TestController) CreateTestUsers(c *fiber.Ctx) error {
	db := config.GetDB()

	// Create test users if they don't exist
	users := []model.Users{
		{
			Name:     "Test User 1",
			Email:    "test1@example.com",
			Phone:    "1234567890",
			Password: "password123",
		},
		{
			Name:     "Test User 2",
			Email:    "test2@example.com",
			Phone:    "1234567891",
			Password: "password123",
		},
	}

	var createdUsers []model.Users

	for _, user := range users {
		// Check if user already exists
		var existingUser model.Users
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			createdUsers = append(createdUsers, existingUser)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return helper.Message500("Failed to hash password")
		}
		user.Password = string(hashedPassword)

		// Create user
		if err := db.Create(&user).Error; err != nil {
			return helper.Message500("Failed to create user: " + err.Error())
		}

		createdUsers = append(createdUsers, user)
	}

	return helper.Message200(c, fiber.Map{
		"users": createdUsers,
	}, "Test users created successfully")
}

// CreateTestChat creates a test chat between two users
func (ctrl *TestController) CreateTestChat(c *fiber.Ctx) error {
	user1IDStr := c.Query("user1_id", "1")
	user2IDStr := c.Query("user2_id", "2")

	user1ID, err := strconv.ParseUint(user1IDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user1_id")
	}

	user2ID, err := strconv.ParseUint(user2IDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user2_id")
	}

	chat, err := ctrl.ChatService.GetOrCreateChat(uint(user1ID), uint(user2ID))
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, chat, "Test chat created successfully")
}

// GetTestToken generates a test JWT token for a user
func (ctrl *TestController) GetTestToken(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return helper.Message400("User ID is required")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user ID")
	}

	// Get user from database
	db := config.GetDB()
	var user model.Users
	if err := db.First(&user, userID).Error; err != nil {
		return helper.Message404("User not found")
	}

	// Generate token
	token, err := ctrl.AuthService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		return helper.Message500("Failed to generate token")
	}

	return helper.Message200(c, fiber.Map{
		"user":  user,
		"token": token,
	}, "Test token generated successfully")
}

// ListTestData shows test users and chats for debugging
func (ctrl *TestController) ListTestData(c *fiber.Ctx) error {
	db := config.GetDB()

	var users []model.Users
	db.Find(&users)

	var chats []model.Chat
	db.Preload("User1").Preload("User2").Find(&chats)

	return helper.Message200(c, fiber.Map{
		"users": users,
		"chats": chats,
	}, "Test data retrieved successfully")
}
