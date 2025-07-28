package model

import "time"

type SocialAuth struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       Users     `json:"user" gorm:"foreignKey:UserID"`
	Provider   string    `json:"provider" gorm:"not null;index:idx_provider_id"`
	ProviderID string    `json:"provider_id" gorm:"not null;index:idx_provider_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SocialAuth) TableName() string {
	return "social_auth"
}
