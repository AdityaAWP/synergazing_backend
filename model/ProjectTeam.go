package model

import "time"

type ProjectRole struct {
	ID             uint                `json:"id" gorm:"primaryKey"`
	ProjectID      uint                `json:"project_id" gorm:"not null"`
	Name           string              `json:"name" gorm:"not null"`
	SlotsAvailable int                 `json:"slots_available" gorm:"not null"`
	Description    string              `json:"description" gorm:"type:text"`
	RequiredSkills []*ProjectRoleSkill `json:"required_skills" gorm:"foreignKey:ProjectRoleID"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func (ProjectRole) TableName() string {
	return "project_roles"
}

type ProjectMember struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	ProjectID     uint        `json:"project_id" gorm:"not null"`
	UserID        uint        `json:"user_id" gorm:"not null"`
	ProjectRoleID uint        `json:"project_role_id" gorm:"not null"`
	Status        string      `json:"status" gorm:"not null;default:'invited'"`
	User          Users       `json:"user" gorm:"foreignKey:UserID"`
	ProjectRole   ProjectRole `json:"project_role" gorm:"foreignKey:ProjectRoleID"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}

// ProjectRoleSkill is the explicit join table for a Role's required skills.
type ProjectRoleSkill struct {
	ProjectRoleID uint  `json:"project_role_id" gorm:"primaryKey"`
	SkillID       uint  `json:"skill_id" gorm:"primaryKey"`
	Skill         Skill `json:"skill" gorm:"foreignKey:SkillID"`
}

func (ProjectRoleSkill) TableName() string {
	return "project_role_skills"
}
