package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	// Nickname  string    `gorm:"not null;default:null;primaryKey"`
	ID        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Chat struct {
	// ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	ID            uuid.UUID
	User1ID       uuid.UUID
	User2ID       uuid.UUID
	User1         User
	User2         User
	LastMessageAt time.Time
	LastMessage   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type Message struct {
	// ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	ID        uuid.UUID
	Content   string
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// type User struct {
// 	gorm.Model
// 	Nickname string
// 	Username string
// 	Password string
// }

// type Message struct {
// 	gorm.Model
// 	Content    string
// 	SenderID   uint
// 	ReceiverID uint
// }

// type Chat struct {
// 	ContactNickname string
// 	ContactID       string
// 	LastMessage     string
// 	LastMessageAt   time.Time
// }
