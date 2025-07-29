package model

import "time"

type Profiles struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	UserID         uint   `json:"user_id" gorm:"not null"`
	ProfilePicture string `json:"profile_picture" gorm:"type:text"`
	User           Users  `json:"user" gorm:"foreignKey:UserID"`

	AboutMe   string `json:"about_me" gorm:"type:text"`
	Location  string `json:"location" gorm:"type:text"`
	Interests string `json:"interests" gorm:"type:text"`
	Academic  string `json:"academic" gorm:"type:text"`

	WebsiteURL    string `json:"website_url" gorm:"type:text"`
	GithubURL     string `json:"github_url" gorm:"type:text"`
	LinkedInURL   string `json:"linkedin_url" gorm:"type:text"`
	InstagramURL  string `json:"instagram_url" gorm:"type:text"`
	PortofolioURL string `json:"portfolio_url" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Profiles) TableName() string {
	return "profiles"
}
