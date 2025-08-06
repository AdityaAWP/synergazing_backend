package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/model"
)

type SkillService struct {
	DB *gorm.DB
}

func NewSkillService(db *gorm.DB) *SkillService {
	return &SkillService{DB: db}
}

func NewSkillServiceDefault() *SkillService {
	db := config.GetDB()
	return &SkillService{DB: db}
}

func (s *SkillService) FindOrCreate(name string) (*model.Skill, error) {
	var skill model.Skill
	if err := s.DB.Where("LOWER(name) = LOWER(?)", name).First(&skill).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newSkill := model.Skill{Name: name}
			if err := s.DB.Create(&newSkill).Error; err != nil {
				return nil, err
			}
			return &newSkill, nil
		}
		return nil, err
	}
	return &skill, nil
}

func (s *SkillService) FindOrCreateWithTx(tx *gorm.DB, name string) (*model.Skill, error) {
	var skill model.Skill
	if err := tx.Where("LOWER(name) = LOWER(?)", name).First(&skill).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newSkill := model.Skill{Name: name}
			if err := tx.Create(&newSkill).Error; err != nil {
				return nil, err
			}
			return &newSkill, nil
		}
		return nil, err
	}
	return &skill, nil
}

func (s *SkillService) GetAllSkills() ([]*model.Skill, error) {
	var skills []*model.Skill
	if err := s.DB.Order("name ASC").Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (s *SkillService) DeleteUserSkills(userId uint, skillName string) error {
	tx := s.DB.Begin()

	if tx.Error != nil {
		return errors.New("Failed to start")
	}

	var skill model.Skill
	if err := tx.Where("LOWER(name) = LOWER(?)", skillName).First(&skill).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("skill '%s' not found", skillName)
		}
		return errors.New("Failed to find skill")
	}

	result := tx.Where("user_id = ? AND skill_id = ?", userId, skill.ID).Delete(&model.UserSkill{})
	if result.Error != nil {
		tx.Rollback()
		return errors.New("Failed to delete user skill")
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("user does not have this skill")
	}
	return tx.Commit().Error
}

func (s *SkillService) GetUserSkills(userId uint) (*model.Users, error) {
	var user model.Users
	if err := s.DB.Preload("UserSkills.Skill").First(&user, userId).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	user.Password = ""
	return &user, nil
}

func (s *SkillService) UpdateUserSkills(userId uint, skillNames []string, proficiencies []int) error {
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

		skill, err := s.FindOrCreateWithTx(tx, skillName)
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
