package model

import "time"

type Users struct {
	ID                  uint         `json:"id" gorm:"primaryKey"`
	Name                string       `json:"name" gorm:"not null"`
	Email               string       `json:"email" gorm:"uniqueIndex;not null"`
	Password            string       `json:"-"`
	Phone               string       `json:"phone" gorm:"uniqueIndex;not null"`
	Role                []*Role      `json:"-" gorm:"many2many:user_roles"`
	StatusCollaboration string       `json:"status_collaboration" gorm:"type:varchar(20);default:'not ready';check:status_collaboration IN ('not ready','ready')"`
	UserSkills          []*UserSkill `json:"user_skills,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`

	PasswordResetToken string
	PasswordResetAt    time.Time
}

func (Users) TableName() string {
	return "users"
}
