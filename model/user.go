package model

import "time"

type Users struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" gorm:"uniqueIndex;not null"`
	Role      []*Role   `json:"-" gorm:"many2many:user_roles"`
	Skills    []*Skill  `json:"skills,omitempty" gorm:"many2many:user_skills;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Users) TableName() string {
	return "users"
}
