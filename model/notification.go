package model

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProjectID *uint     `json:"project_id,omitempty"`
	Type      string    `json:"type" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	Data      string    `json:"data,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User    Users    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

func (Notification) TableName() string {
	return "notifications"
}

// Notification types constants
const (
	NotificationTypeDeadlineApproaching = "deadline_approaching"
	NotificationTypeUserRegistered      = "user_registered"
	NotificationTypeUserAccepted        = "user_accepted"
	NotificationTypeUserRejected        = "user_rejected"
	NotificationTypeProjectStatusChange = "project_status_change"
	NotificationTypeProjectUpdated      = "project_updated"
	NotificationTypeTeamMemberLeft      = "team_member_left"
	NotificationTypeProjectCompleted    = "project_completed"
	NotificationTypeRoleAssigned        = "role_assigned"
	NotificationTypeInvitationReceived  = "invitation_received"
)
