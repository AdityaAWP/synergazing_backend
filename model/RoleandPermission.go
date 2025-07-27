package model

import (
	"time"
)

type Role struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name" gorm:"unique;not null;size:50"`
	Description string        `json:"description" gorm:"size:255"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []*Users      `json:"users" gorm:"many2many:user_roles;"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null;size:100"`
	Slug      string    `json:"slug" gorm:"unique;not null;size:100"`
	Group     string    `json:"group" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}
