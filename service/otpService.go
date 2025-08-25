package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/model"
)

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (s *OTPService) GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}

func (s *OTPService) SendOTP(email, purpose string) error {
	db := config.GetDB()

	code, err := s.GenerateOTP()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	db.Where("email = ? AND purpose = ? AND is_used = ?", email, purpose, false).Delete(&model.OTP{})

	otp := model.OTP{
		Email:     email,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(time.Minute * 3),
		IsUsed:    false,
	}

	if err := db.Create(&otp).Error; err != nil {
		log.Printf("Database error creating OTP: %v", err)
		return errors.New("failed to create OTP")
	}

	go s.sendOTPEmail(email, code, purpose)

	return nil
}

func (s *OTPService) VerifyOTP(email, code, purpose string) error {
	db := config.GetDB()

	var otp model.OTP
	if err := db.Where("email = ? AND code = ? AND purpose = ? AND is_used = ?",
		email, code, purpose, false).First(&otp).Error; err != nil {
		return errors.New("invalid OTP code")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	otp.IsUsed = true
	if err := db.Save(&otp).Error; err != nil {
		log.Printf("Database error marking OTP as used: %v", err)
		return errors.New("failed to verify OTP")
	}

	return nil
}

func (s *OTPService) CleanupExpiredOTPs() {
	db := config.GetDB()

	result := db.Where("expires_at < ?", time.Now()).Delete(&model.OTP{})
	if result.Error != nil {
		log.Printf("Error cleaning up expired OTPs: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired OTP records", result.RowsAffected)
	}
}

func (s *OTPService) sendOTPEmail(email, code, purpose string) {
	emailHost := os.Getenv("EMAIL_HOST")
	emailPortStr := os.Getenv("EMAIL_PORT")
	emailUser := os.Getenv("EMAIL_USERNAME")
	emailPass := os.Getenv("EMAIL_PASSWORD")

	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Printf("Error: Could not parse EMAIL_PORT from .env file: %v", err)
		return
	}

	m := gomail.NewMessage()

	m.SetHeader("From", fmt.Sprintf("Synergazing <%s>", emailUser))
	m.SetHeader("To", email)

	var subject, htmlBody, plainBody string

	switch purpose {
	case "registration":
		subject = "Email Verification - Complete Your Registration"
		htmlBody = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Welcome to Synergazing!</h2>
			<p>Hi there,</p>
			<p>Thank you for registering with Synergazing. To complete your registration, please verify your email address using the verification code below:</p>
			<div style="background-color: #f8f9fa; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px;">
				<h1 style="color: #007bff; font-size: 32px; letter-spacing: 8px; margin: 0;">%s</h1>
			</div>
			<p>This verification code will expire in <strong>3 minutes</strong>.</p>
			<p>If you did not create an account with Synergazing, please ignore this email.</p>
			<br>
			<p>Best regards,</p>
			<p>The Synergazing Team</p>
		</div>
		`, code)

		plainBody = fmt.Sprintf(
			"Welcome to Synergazing!\n\n"+
				"Thank you for registering with Synergazing. To complete your registration, please verify your email address using the verification code below:\n\n"+
				"Verification Code: %s\n\n"+
				"This verification code will expire in 3 minutes.\n\n"+
				"If you did not create an account with Synergazing, please ignore this email.\n\n"+
				"Best regards,\nThe Synergazing Team", code)

	case "password_reset":
		subject = "Password Reset Verification Code"
		htmlBody = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Password Reset Request</h2>
			<p>Hi,</p>
			<p>We received a request to reset your password. Please use the verification code below to proceed:</p>
			<div style="background-color: #f8f9fa; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px;">
				<h1 style="color: #dc3545; font-size: 32px; letter-spacing: 8px; margin: 0;">%s</h1>
			</div>
			<p>This verification code will expire in <strong>3 minutes</strong>.</p>
			<p>If you did not request a password reset, please ignore this email and your password will remain unchanged.</p>
			<br>
			<p>Best regards,</p>
			<p>The Synergazing Team</p>
		</div>
		`, code)

		plainBody = fmt.Sprintf(
			"Password Reset Request\n\n"+
				"We received a request to reset your password. Please use the verification code below to proceed:\n\n"+
				"Verification Code: %s\n\n"+
				"This verification code will expire in 3 minutes.\n\n"+
				"If you did not request a password reset, please ignore this email and your password will remain unchanged.\n\n"+
				"Best regards,\nThe Synergazing Team", code)

	default:
		subject = "Verification Code"
		htmlBody = fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.6; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Verification Required</h2>
			<p>Hi,</p>
			<p>Please use the verification code below:</p>
			<div style="background-color: #f8f9fa; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px;">
				<h1 style="color: #007bff; font-size: 32px; letter-spacing: 8px; margin: 0;">%s</h1>
			</div>
			<p>This verification code will expire in <strong>3 minutes</strong>.</p>
			<br>
			<p>Best regards,</p>
			<p>The Synergazing Team</p>
		</div>
		`, code)

		plainBody = fmt.Sprintf(
			"Verification Required\n\n"+
				"Please use the verification code below:\n\n"+
				"Verification Code: %s\n\n"+
				"This verification code will expire in 3 minutes.\n\n"+
				"Best regards,\nThe Synergazing Team", code)
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	m.AddAlternative("text/plain", plainBody)

	d := gomail.NewDialer(emailHost, emailPort, emailUser, emailPass)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send OTP email to %s: %v", email, err)
	} else {
		log.Printf("OTP email sent successfully to %s for purpose: %s", email, purpose)
	}
}
