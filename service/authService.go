package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

type AuthService struct {
	OTPService *OTPService
}

func NewAuthService(otpService *OTPService) *AuthService {
	return &AuthService{
		OTPService: otpService,
	}
}

// InitiateRegistration validates user data and sends OTP for email verification
func (s *AuthService) InitiateRegistration(name, email, password, phone string) error {
	if name == "" {
		return errors.New("Name is required")
	}
	if email == "" {
		return errors.New("Email is required")
	}
	if password == "" {
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}
	if phone == "" {
		return errors.New("Phone number is required")
	}

	db := config.GetDB()

	// Check if user already exists
	var existUser model.Users
	if err := db.Where("email = ?", email).First(&existUser).Error; err == nil {
		return errors.New("Email already exists")
	}

	// Send OTP for email verification
	return s.OTPService.SendOTP(email, "registration")
}

// CompleteRegistration creates user account after OTP verification
func (s *AuthService) CompleteRegistration(name, email, password, phone, otpCode string) (*model.Users, error) {
	// Verify OTP first
	if err := s.OTPService.VerifyOTP(email, otpCode, "registration"); err != nil {
		return nil, err
	}

	db := config.GetDB()

	// Double-check that user doesn't exist (in case they were created between initiate and complete)
	var existUser model.Users
	if err := db.Where("email = ?", email).First(&existUser).Error; err == nil {
		return nil, errors.New("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}

	user := model.Users{
		Name:            name,
		Email:           email,
		Password:        string(hashedPassword),
		Phone:           phone,
		IsEmailVerified: true, // Mark as verified since OTP was successful
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Database error creating user: %v", err)
		return nil, fmt.Errorf("Failed to create user: %v", err)
	}

	user.Password = ""
	return &user, nil
}

// Register (legacy method for backward compatibility)
func (s *AuthService) Register(name, email, password string, phone string) (*model.Users, error) {
	if name == "" {
		return nil, errors.New("Name is required")
	}
	if email == "" {
		return nil, errors.New("Email is required")
	}
	if password == "" {
		return nil, errors.New("Password is required")
	}
	if len(password) < 8 {
		return nil, errors.New("Password must be at least 8 characters")
	}
	if phone == "" {
		return nil, errors.New("Phone number is required")
	}

	db := config.GetDB()

	var existUser model.Users
	if err := db.Where("email = ?", email).First(&existUser).Error; err == nil {
		return nil, errors.New("Email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Failed to Hash password")
	}

	user := model.Users{
		Name:            name,
		Email:           email,
		Password:        string(hashedPassword),
		Phone:           phone,
		IsEmailVerified: true, // Direct registration doesn't require email verification
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Database error creating user: %v", err)
		return nil, fmt.Errorf("Failed to create user: %v", err)
	}
	user.Password = ""
	return &user, nil
}

func (s *AuthService) Login(email, password string) (string, *model.Users, error) {
	db := config.GetDB()

	var user model.Users
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, errors.New("Invalid Credential Email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("Invalid Credential Password")
	}

	token, err := helper.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	user.Password = ""
	return token, &user, nil
}

func (s *AuthService) Logout(token string) error {
	return nil
}

// ResendOTP resends OTP for the given email and purpose
func (s *AuthService) ResendOTP(email, purpose string) error {
	return s.OTPService.SendOTP(email, purpose)
}

// VerifyEmailWithOTP verifies email using OTP code
func (s *AuthService) VerifyEmailWithOTP(email, otpCode string) error {
	if err := s.OTPService.VerifyOTP(email, otpCode, "registration"); err != nil {
		return err
	}

	db := config.GetDB()
	var user model.Users
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("User not found")
	}

	user.IsEmailVerified = true
	if err := db.Save(&user).Error; err != nil {
		return errors.New("Failed to update email verification status")
	}

	return nil
}

func (s *AuthService) GenerateTokenForUser(userID uint, email string) (string, error) {
	token, err := helper.GenerateJWTToken(userID, email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	db := config.GetDB()
	var user model.Users
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Password reset requested for non-existent email: %s", email)
		return nil
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating token: %v", err)
		return errors.New("could not generate token")
	}
	token := hex.EncodeToString(b)

	user.PasswordResetToken = token
	user.PasswordResetAt = time.Now().Add(time.Minute * 5)
	if err := db.Save(&user).Error; err != nil {
		log.Printf("Database error saving reset token: %v", err)
		return errors.New("failed to save reset token")
	}

	go helper.SendPasswordResetEmail(user.Email, token)

	return nil
}

func (s *AuthService) ResetPassword(token, password, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("passwords do not match")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	db := config.GetDB()
	var user model.Users
	if err := db.Where("password_reset_token = ?", token).First(&user).Error; err != nil {
		return errors.New("invalid token")
	}

	if time.Now().After(user.PasswordResetAt) {
		return errors.New("token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.PasswordResetToken = ""
	if err := db.Save(&user).Error; err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// GetEmailVerificationStatus returns the email verification status of a user
func (s *AuthService) GetEmailVerificationStatus(email string) (bool, error) {
	db := config.GetDB()
	var user model.Users
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return false, errors.New("User not found")
	}

	return user.IsEmailVerified, nil
}
