package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type SocialAuthService struct {
	db *gorm.DB
}

func NewSocialAuthService(db *gorm.DB) *SocialAuthService {
	return &SocialAuthService{
		db: db,
	}
}

func (s *SocialAuthService) HandleProviderCallback(provider, providerID, name, email string) (*model.Users, error) {
	log.Printf("🔄 HandleProviderCallback: provider=%s, providerID=%s, name=%s, email=%s", provider, providerID, name, email)

	var socialAuth model.SocialAuth

	// Check if this social auth already exists
	if err := s.db.Preload("User").Where("provider = ? AND provider_id = ?", provider, providerID).First(&socialAuth).Error; err == nil {
		log.Printf("✅ Found existing social auth for provider_id=%s, returning user_id=%d", providerID, socialAuth.User.ID)
		return &socialAuth.User, nil
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("❌ Database error when checking social auth: %v", err)
		return nil, err
	}

	// Check if user already exists by email
	var user model.Users
	if err := s.db.Where("email = ?", email).First(&user).Error; err == nil {
		log.Printf("✅ Found existing user by email=%s, user_id=%d", email, user.ID)

		// Update existing user to be email verified for OAuth login
		user.IsEmailVerified = true
		if err := s.db.Save(&user).Error; err != nil {
			log.Printf("❌ Failed to update user email verification: %v", err)
			return nil, err
		}

		newSocialAuth := model.SocialAuth{
			UserID:     user.ID,
			Provider:   provider,
			ProviderID: providerID,
		}
		if err := s.db.Create(&newSocialAuth).Error; err != nil {
			log.Printf("❌ Failed to create social auth for existing user: %v", err)
			return nil, err
		}
		log.Printf("✅ Successfully linked social auth to existing user_id=%d", user.ID)
		return &user, nil
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("❌ Database error when checking user by email: %v", err)
		return nil, err
	}

	log.Printf("🆕 Creating new user for email=%s", email)

	tx := s.db.Begin()

	// Generate random Indonesian phone number for OAuth users
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(900000000) + 100000000 // 9-digit number (100000000-999999999)
	phoneNumber := fmt.Sprintf("+628%d", randomNumber)

	newUser := model.Users{
		Name:            name,
		Email:           email,
		Phone:           phoneNumber, // Use random Indonesian phone number
		IsEmailVerified: true,        // OAuth users are automatically email verified
	}

	if err := tx.Create(&newUser).Error; err != nil {
		log.Printf("❌ Failed to create new user: %v", err)
		tx.Rollback()
		return nil, err
	}
	log.Printf("✅ Created new user_id=%d", newUser.ID)

	newSocialAuth := model.SocialAuth{
		UserID:     newUser.ID,
		Provider:   provider,
		ProviderID: providerID,
	}
	if err := tx.Create(&newSocialAuth).Error; err != nil {
		log.Printf("❌ Failed to create social auth for new user: %v", err)
		tx.Rollback()
		return nil, err
	}
	log.Printf("✅ Created social auth for user_id=%d", newUser.ID)

	if err := tx.Commit().Error; err != nil {
		log.Printf("❌ Failed to commit transaction: %v", err)
		return nil, err
	}

	log.Printf("✅ Successfully created new user and social auth, user_id=%d", newUser.ID)
	return &newUser, nil
}
