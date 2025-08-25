package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type AuthController struct {
	AuthService *service.AuthService
	OTPService  *service.OTPService
}

func NewAuthController(AuthService *service.AuthService, OTPService *service.OTPService) *AuthController {
	return &AuthController{
		AuthService: AuthService,
		OTPService:  OTPService,
	}
}

// InitiateRegistration starts the registration process and sends OTP
func (ctrl *AuthController) InitiateRegistration(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	phone := c.FormValue("phone")

	if name == "" || email == "" || password == "" || phone == "" {
		return helper.Message400("All fields are required")
	}

	err := ctrl.AuthService.InitiateRegistration(name, email, password, phone)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "OTP sent to your email. Please verify to complete registration.")
}

// CompleteRegistration completes the registration process after OTP verification
func (ctrl *AuthController) CompleteRegistration(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	phone := c.FormValue("phone")
	otpCode := c.FormValue("otp_code")

	if name == "" || email == "" || password == "" || phone == "" || otpCode == "" {
		return helper.Message400("All fields are required")
	}

	user, err := ctrl.AuthService.CompleteRegistration(name, email, password, phone, otpCode)
	if err != nil {
		return helper.Message400(err.Error())
	}

	token, err := ctrl.AuthService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		return helper.Message500("Token generation failed")
	}

	return helper.Message201(c, fiber.Map{
		"user":  user,
		"token": token,
	}, "Registration completed successfully")
}

// Register (legacy method for backward compatibility)
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	phone := c.FormValue("phone")

	if name == "" || email == "" || password == "" || phone == "" {
		return helper.Message400("All fields are required")
	}

	user, err := ctrl.AuthService.Register(name, email, password, phone)
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
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return helper.Message400("Email and password are required")
	}

	token, user, err := ctrl.AuthService.Login(email, password)
	if err != nil {
		return helper.Message401(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"user":  user,
		"token": token,
	}, "Login Successful")
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
		"message": "Logout Successful",
	})
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return helper.Message400("Email is required")
	}

	err := ctrl.AuthService.ForgotPassword(email)
	if err != nil {
		return helper.Message500(err.Error())
	}

	return helper.Message200(c, nil, "If an account with that email exists, a password reset link has been sent.")
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	token := c.FormValue("token")
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("passwordConfirm")

	if token == "" || password == "" || passwordConfirm == "" {
		return helper.Message400("All fields are required")
	}

	err := ctrl.AuthService.ResetPassword(token, password, passwordConfirm)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Password has been reset successfully.")
}

// ResendOTP resends OTP for the given email and purpose
func (ctrl *AuthController) ResendOTP(c *fiber.Ctx) error {
	email := c.FormValue("email")
	purpose := c.FormValue("purpose")

	if email == "" || purpose == "" {
		return helper.Message400("Email and purpose are required")
	}

	// Validate purpose
	if purpose != "registration" && purpose != "password_reset" {
		return helper.Message400("Invalid purpose. Must be 'registration' or 'password_reset'")
	}

	err := ctrl.AuthService.ResendOTP(email, purpose)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "OTP resent successfully")
}

// VerifyOTP verifies OTP code for any purpose
func (ctrl *AuthController) VerifyOTP(c *fiber.Ctx) error {
	email := c.FormValue("email")
	otpCode := c.FormValue("otp_code")
	purpose := c.FormValue("purpose")

	if email == "" || otpCode == "" || purpose == "" {
		return helper.Message400("Email, OTP code, and purpose are required")
	}

	err := ctrl.OTPService.VerifyOTP(email, otpCode, purpose)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "OTP verified successfully")
}
