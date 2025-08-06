package service

import (
	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type SkillService struct {
	db *gorm.DB
}

func NewSkillService(db *gorm.DB) *SkillService {
	return &SkillService{db: db}
}

func (s *SkillService) FindOrCreate(name string) (*model.Skill, error) {
	var skill model.Skill
	if err := s.db.Where("LOWER(name) = LOWER(?)", name).First(&skill).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newSkill := model.Skill{Name: name}
			if err := s.db.Create(&newSkill).Error; err != nil {
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
	if err := s.db.Order("name ASC").Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
