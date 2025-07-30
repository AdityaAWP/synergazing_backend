package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
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
		Phone    int    `json:"phone" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Message400("Invalid request Body")
	}

	user, err := ctrl.AuthService.Register(req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		return helper.Message400(err.Error())
	}

	token, err := ctrl.AuthService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		return helper.Message500("Token generation failed")
	}

	return helper.Message201(c, fiber.Map{
		"users": user,
		"token": token,
	}, "User registered successfully")
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" validate:"required, email"`
		Password string `json:"password" validate:"required, min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Message400("Invalid request body")
	}

	token, user, err := ctrl.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return helper.Message401(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"user":  user,
		"token": token,
	}, "Login Succesful")
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return helper.Message400("No token provided")
	}
	if len(token) > 7 && token[:7] == "Bearer" {
		token = token[:7]
	}

	err := ctrl.AuthService.Logout(token)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Logout Succesfull",
	})
}
