package model

import "time"

type UserSkill struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	SkillID     uint      `json:"skill_id" gorm:"not null;index"`
	Proficiency int       `json:"proficiency" gorm:"not null;check:proficiency >= 0 AND proficiency <= 100"`
	User        Users     `json:"-" gorm:"foreignKey:UserID"`
	Skill       Skill     `json:"skill" gorm:"foreignKey:SkillID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UserSkill) TableName() string {
	return "user_skills"
}
