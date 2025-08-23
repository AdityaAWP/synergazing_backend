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

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

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
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Phone:    phone,
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
		return "", nil, errors.New("Invalid Credetial Email")
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
