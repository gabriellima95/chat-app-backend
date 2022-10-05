package storage

import (
	"log"
	"msn/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericChatRepository struct {
	DB *gorm.DB
}

func NewGenericChatRepository(db *gorm.DB) GenericChatRepository {
	return GenericChatRepository{
		DB: db,
	}
}

func (r GenericChatRepository) Create(chat *models.GenericChat) error {
	chat.ID = uuid.New()
	return r.DB.Create(chat).Error
}

func (r GenericChatRepository) Update(chat *models.GenericChat) error {
	err := r.DB.Save(chat).Error
	if err != nil {
		log.Printf("Failed to run Update error: %v", err)
		return err
	}

	return nil
}

func (r GenericChatRepository) GetByID(id uuid.UUID) (*models.GenericChat, error) {
	var chat models.GenericChat
	err := r.DB.Preload("Users").First(&chat, "id = ?", id).Error
	if err != nil {
		log.Printf("Failed to run GetByID error: %v", err)
		return nil, err
	}
	return &chat, nil
}
