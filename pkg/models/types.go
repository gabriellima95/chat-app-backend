package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	// Nickname  string    `gorm:"not null;default:null;primaryKey"`
	ID        uuid.UUID `gorm:"not null;default:null;primaryKey"`
	Username  string    `gorm:"not null;unique;default:null"`
	Password  string    `gorm:"not null;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Chat struct {
	// ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	ID            uuid.UUID `gorm:"not null;default:null;primaryKey"`
	User1ID       uuid.UUID `gorm:"not null;default:null"`
	User2ID       uuid.UUID `gorm:"not null;default:null"`
	User1         User      `gorm:"foreignKey:User1ID"`
	User2         User      `gorm:"foreignKey:User2ID"`
	LastMessageAt time.Time `gorm:"not null;default:null"`
	LastMessage   string    `gorm:"not null;default:null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type Message struct {
	// ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	ID        uuid.UUID `gorm:"not null;default:null;primaryKey"`
	Content   string    `gorm:"not null;default:null"`
	ChatID    uuid.UUID `gorm:"not null;default:null"`
	SenderID  uuid.UUID `gorm:"not null;default:null"`
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