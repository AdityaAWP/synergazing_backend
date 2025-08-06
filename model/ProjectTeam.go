package model

import "time"

type ProjectRole struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ProjectID      uint      `json:"project_id" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	SlotsAvailable int       `json:"slots_available" gorm:"not null"`
	Description    string    `json:"description" gorm:"type:text"`
	RequiredSkills []*Skill  `json:"required_skills" gorm:"many2many:project_role_skills;"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ProjectRole) TableName() string {
	return "project_roles"
}

type ProjectMember struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	ProjectID     uint   `json:"project_id" gorm:"not null"`
	UserID        uint   `json:"user_id" gorm:"not null"`
	ProjectRoleID uint   `json:"project_role_id" gorm:"not null"`
	Status        string `json:"status" gorm:"not null;default:'invited'"`

	User        Users       `json:"user" gorm:"foreignKey:UserID"`
	ProjectRole ProjectRole `json:"project_role" gorm:"foreignKey:ProjectRoleID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}
