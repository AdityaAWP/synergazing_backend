package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/service"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(AuthService *service.AuthService) *AuthController {
	return &AuthController{
		AuthService: AuthService,
	}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required, email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request Body",
		})
	}

	user, err := ctrl.AuthService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := ctrl.AuthService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Token generation failed",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" validate:"required, email"`
		Password string `json:"password" validate:"required, min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, user, err := ctrl.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Login Succesful",
		"user":    user,
		"token":   token,
	})
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "No token provided",
		})
	}
	if len(token) > 7 && token[:7] == "Bearer" {
		token = token[:7]
	}

	err := ctrl.AuthService.Logout(token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logout Succesfull",
	})
}
