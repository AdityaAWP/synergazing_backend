package model

import "time"

type Chat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	User1ID   uint      `json:"user1_id" gorm:"not null"`
	User2ID   uint      `json:"user2_id" gorm:"not null"`
	User1     Users     `json:"user1" gorm:"foreignKey:User1ID"`
	User2     Users     `json:"user2" gorm:"foreignKey:User2ID"`
	Messages  []Message `json:"messages,omitempty" gorm:"foreignKey:ChatID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Chat) TableName() string {
	return "chats"
}

type Message struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ChatID    uint      `json:"chat_id" gorm:"not null"`
	SenderID  uint      `json:"sender_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	Chat      Chat      `json:"chat" gorm:"foreignKey:ChatID"`
	Sender    Users     `json:"sender" gorm:"foreignKey:SenderID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Message) TableName() string {
	return "messages"
}
