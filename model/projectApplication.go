package model

import "time"

type ProjectApplication struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	ProjectID     uint   `json:"project_id" gorm:"not null"`
	UserID        uint   `json:"user_id" gorm:"not null"`
	ProjectRoleID uint   `json:"project_role_id" gorm:"not null"`
	Status        string `json:"status" gorm:"not null;default:'pending'"`

	// Enhanced application fields
	WhyInterested    string `json:"why_interested" gorm:"type:text;not null"`
	SkillsExperience string `json:"skills_experience" gorm:"type:text;not null"`
	Contribution     string `json:"contribution" gorm:"type:text;not null"`

	AppliedAt   time.Time  `json:"applied_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  *uint      `json:"reviewed_by,omitempty"`
	ReviewNotes string     `json:"review_notes" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	Project     Project     `json:"project" gorm:"foreignKey:ProjectID"`
	User        Users       `json:"user" gorm:"foreignKey:UserID"`
	ProjectRole ProjectRole `json:"project_role" gorm:"foreignKey:ProjectRoleID"`
	Reviewer    *Users      `json:"reviewer,omitempty" gorm:"foreignKey:ReviewedBy"`
}

func (ProjectApplication) TableName() string {
	return "project_applications"
}

// Application status constants
const (
	ApplicationStatusPending   = "pending"
	ApplicationStatusAccepted  = "accepted"
	ApplicationStatusRejected  = "rejected"
	ApplicationStatusWithdrawn = "withdrawn"
)
