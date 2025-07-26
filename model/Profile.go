package model

import "time"

type Profiles struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:text"`
	User           Users     `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Profiles) TableName() string {
	return "profiles"
}
