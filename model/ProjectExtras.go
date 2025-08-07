package model

import "time"

type ProjectCondition struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	ProjectID   uint   `json:"project_id" gorm:"not null"`
	Description string `json:"description" gorm:"type:text;not null"`
}

func (ProjectCondition) TableName() string {
	return "project_conditions"
}

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}

type ProjectRequiredSkill struct {
	ProjectID uint  `json:"project_id" gorm:"primaryKey"`
	SkillID   uint  `json:"skill_id" gorm:"primaryKey"`
	Skill     Skill `json:"skill" gorm:"foreignKey:SkillID"`
}

func (ProjectRequiredSkill) TableName() string {
	return "project_required_skills"
}

type ProjectTag struct {
	ProjectID uint `json:"project_id" gorm:"primaryKey"`
	TagID     uint `json:"tag_id" gorm:"primaryKey"`
	Tag       Tag  `json:"tag" gorm:"foreignKey:TagID"`
}

func (ProjectTag) TableName() string {
	return "project_tags"
}
