package service

import (
	"fmt"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

type UpdateProfileDTO struct {
	Name           *string
	Email          *string
	Phone          *string
	Password       *string
	AboutMe        *string
	Location       *string
	Interests      *string
	Academic       *string
	WebsiteURL     *string
	GithubURL      *string
	LinkedInURL    *string
	InstagramURL   *string
	PortfolioURL   *string
	ProfilePicture *multipart.FileHeader
	CVFile         *multipart.FileHeader
}

type ProfileService struct {
	DB           *gorm.DB
	SkillService *SkillService
}

func NewProfileService() *ProfileService {
	db := config.GetDB()
	return &ProfileService{
		DB:           db,
		SkillService: NewSkillService(db),
	}
}

func (s *ProfileService) GetUserProfile(userId uint) (*model.Users, *model.Profiles, error) {
	var user model.Users
	var profile model.Profiles

	if err := s.DB.Preload("Role").Preload("UserSkills.Skill").First(&user, userId).Error; err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	if err := s.DB.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			profile = model.Profiles{UserID: user.ID}
			if err := s.DB.Create(&profile).Error; err != nil {
				return nil, nil, fmt.Errorf("failed to create profile")
			}
		} else {
			return nil, nil, err
		}
	}
	user.Password = ""
	return &user, &profile, nil
}

func (s *ProfileService) UpdateUserProfile(userId uint, data *UpdateProfileDTO) (*model.Users, *model.Profiles, error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.Users
	if err := tx.First(&user, userId).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("user not found")
	}

	var profile model.Profiles
	if err := tx.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			profile = model.Profiles{UserID: userId}
			if err := tx.Create(&profile).Error; err != nil {
				tx.Rollback()
				return nil, nil, fmt.Errorf("failed to create profile")
			}
		} else {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if data.Name != nil {
		user.Name = *data.Name
	}
	if data.Phone != nil {
		user.Phone = *data.Phone
	}
	if data.Email != nil {
		var existingUser model.Users
		if err := tx.Where("email = ? AND id != ?", *data.Email, userId).First(&existingUser).Error; err == nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("email already exists")
		}
		user.Email = *data.Email
	}
	if data.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}

	if data.AboutMe != nil {
		profile.AboutMe = *data.AboutMe
	}
	if data.Location != nil {
		profile.Location = *data.Location
	}
	if data.Interests != nil {
		profile.Interests = *data.Interests
	}
	if data.Academic != nil {
		profile.Academic = *data.Academic
	}
	if data.WebsiteURL != nil {
		profile.WebsiteURL = *data.WebsiteURL
	}
	if data.GithubURL != nil {
		profile.GithubURL = *data.GithubURL
	}
	if data.LinkedInURL != nil {
		profile.LinkedInURL = *data.LinkedInURL
	}
	if data.InstagramURL != nil {
		profile.InstagramURL = *data.InstagramURL
	}
	if data.PortfolioURL != nil {
		profile.PortfolioURL = *data.PortfolioURL
	}

	var newProfilePicPath string
	if data.ProfilePicture != nil {
		if profile.ProfilePicture != "" {
			helper.DeleteFile(profile.ProfilePicture)
		}
		filePath, err := helper.UploadFile(data.ProfilePicture, "profile")
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
		newProfilePicPath = filePath
		profile.ProfilePicture = newProfilePicPath
	}

	var newCvPath string
	if data.CVFile != nil {
		if profile.CVFile != "" {
			helper.DeleteFile(profile.CVFile)
		}
		filePath, err := helper.UploadFile(data.CVFile, "cv")
		if err != nil {
			tx.Rollback()
			if newProfilePicPath != "" {
				helper.DeleteFile(newProfilePicPath)
			}
			return nil, nil, err
		}
		newCvPath = filePath
		profile.CVFile = newCvPath
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		if newProfilePicPath != "" {
			helper.DeleteFile(newProfilePicPath)
		}
		if newCvPath != "" {
			helper.DeleteFile(newCvPath)
		}
		return nil, nil, err
	}
	if err := tx.Save(&profile).Error; err != nil {
		tx.Rollback()
		if newProfilePicPath != "" {
			helper.DeleteFile(newProfilePicPath)
		}
		if newCvPath != "" {
			helper.DeleteFile(newCvPath)
		}
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		if newProfilePicPath != "" {
			helper.DeleteFile(newProfilePicPath)
		}
		if newCvPath != "" {
			helper.DeleteFile(newCvPath)
		}
		return nil, nil, err
	}

	// Reload the profile to get updated data with proper relationships
	if err := s.DB.Preload("User").Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to reload profile")
	}

	user.Password = ""
	profile.User.Password = ""
	return &user, &profile, nil
}

func (s *ProfileService) GetPublicUserProfile(userId uint) (*model.Users, *model.Profiles, error) {
	var user model.Users
	if err := s.DB.Preload("UserSkills.Skill").First(&user, userId).Error; err != nil {
		return nil, nil, err
	}

	var profile model.Profiles
	if err := s.DB.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user, &model.Profiles{UserID: userId}, nil
		}
		return nil, nil, err
	}

	user.Password = ""
	return &user, &profile, nil
}

func (s *ProfileService) GetCVFilePath(userId uint) (string, error) {
	var profile model.Profiles
	if err := s.DB.Select("cv_file").Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return "", err
	}
	return profile.CVFile, nil
}

func (s *ProfileService) DeleteProfilePicture(userId uint) error {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var profile model.Profiles
	if err := tx.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("profile not found")
	}

	if profile.ProfilePicture == "" {
		tx.Rollback()
		return fmt.Errorf("no profile picture to delete")
	}

	if err := helper.DeleteFile(profile.ProfilePicture); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	if err := tx.Model(&profile).Update("profile_picture", "").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update profile record: %w", err)
	}

	return tx.Commit().Error
}

func (s *ProfileService) UpdateUserSkills(userId uint, skillNames []string, proficiencies []int) error {
	if len(skillNames) != len(proficiencies) {
		return fmt.Errorf("skill names and proficiencies length mismatch")
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", userId).Delete(&model.UserSkill{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing skills: %w", err)
	}

	for i, skillName := range skillNames {
		if skillName == "" {
			continue
		}

		skill, err := s.SkillService.FindOrCreateWithTx(tx, skillName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find or create skill '%s': %w", skillName, err)
		}

		userSkill := model.UserSkill{
			UserID:      userId,
			SkillID:     skill.ID,
			Proficiency: proficiencies[i],
		}

		if err := tx.Create(&userSkill).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create user skill association: %w", err)
		}
	}

	return tx.Commit().Error
}

func (s *ProfileService) DeleteCVFile(userId uint) error {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var profile model.Profiles
	if err := tx.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("profile not found")
	}

	if profile.CVFile == "" {
		tx.Rollback()
		return fmt.Errorf("no CV file to delete")
	}

	if err := helper.DeleteFile(profile.CVFile); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	if err := tx.Model(&profile).Update("cv_file", "").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update profile record: %w", err)
	}

	return tx.Commit().Error
}
