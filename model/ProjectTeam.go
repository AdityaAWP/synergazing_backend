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
	ID              uint                  `json:"id" gorm:"primaryKey"`
	ProjectID       uint                  `json:"project_id" gorm:"not null"`
	UserID          uint                  `json:"user_id" gorm:"not null"`
	ProjectRoleID   uint                  `json:"project_role_id" gorm:"not null"`
	Status          string                `json:"status" gorm:"not null;default:'invited'"`
	RoleDescription string                `json:"role_description" gorm:"type:text"`
	Project         Project               `json:"project" gorm:"foreignKey:ProjectID"`
	User            Users                 `json:"user" gorm:"foreignKey:UserID"`
	ProjectRole     ProjectRole           `json:"project_role" gorm:"foreignKey:ProjectRoleID"`
	MemberSkills    []*ProjectMemberSkill `json:"member_skills" gorm:"foreignKey:ProjectMemberID"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}

type ProjectRoleSkill struct {
	ProjectRoleID uint  `json:"project_role_id" gorm:"primaryKey"`
	SkillID       uint  `json:"skill_id" gorm:"primaryKey"`
	Skill         Skill `json:"skill" gorm:"foreignKey:SkillID"`
}

func (ProjectRoleSkill) TableName() string {
	return "project_role_skills"
}

type ProjectMemberSkill struct {
	ProjectMemberID uint  `json:"project_member_id" gorm:"primaryKey"`
	SkillID         uint  `json:"skill_id" gorm:"primaryKey"`
	Skill           Skill `json:"skill" gorm:"foreignKey:SkillID"`
}

func (ProjectMemberSkill) TableName() string {
	return "project_member_skills"
}
