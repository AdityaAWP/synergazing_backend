package service

import (
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
	var socialAuth model.SocialAuth

	if err := s.db.Preload("User").Where("provider = ? AND provider_id = ?", provider, providerID).First(&socialAuth).Error; err == nil {
		return &socialAuth.User, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var user model.Users
	if err := s.db.Where("email = ?", email).First(&user).Error; err == nil {
		newSocialAuth := model.SocialAuth{
			UserID:     user.ID,
			Provider:   provider,
			ProviderID: providerID,
		}
		if err := s.db.Create(&newSocialAuth).Error; err != nil {
			return nil, err
		}
		return &user, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tx := s.db.Begin()

	newUser := model.Users{
		Name:  name,
		Email: email,
	}
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	newSocialAuth := model.SocialAuth{
		UserID:     newUser.ID,
		Provider:   provider,
		ProviderID: providerID,
	}
	if err := tx.Create(&newSocialAuth).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}
