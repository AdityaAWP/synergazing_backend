package model

import "time"

type OTP struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null;index"`
	Code      string    `json:"code" gorm:"not null"`
	Purpose   string    `json:"purpose" gorm:"not null;default:'registration';check:purpose IN ('registration','password_reset','email_change')"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (OTP) TableName() string {
	return "otps"
}
