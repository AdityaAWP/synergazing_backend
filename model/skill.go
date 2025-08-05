package model

import "time"

type Skill struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null;size:100"`

	UserSkills []*UserSkill `json:"user_skills,omitempty" gorm:"foreignKey:SkillID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Skill) TableName() string {
	return "skill"
}
