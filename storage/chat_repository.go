package storage

import (
	"log"
	"msn/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	DB *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return ChatRepository{
		DB: db,
	}
}

func (r ChatRepository) Create(chat *models.Chat) error {
	chat.ID = uuid.New()
	return r.DB.Create(chat).Error
}

func (r ChatRepository) ListByUserID(id uuid.UUID) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.DB.Preload("User1").Preload("User2").Where("user1_id = ? OR user2_id = ?", id.String(), id.String()).Order("last_message_at desc").Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r ChatRepository) Update(chat *models.Chat) error {
	err := r.DB.Save(chat).Error
	if err != nil {
		log.Printf("Failed to run Update error: %v", err)
		return err
	}

	return nil
}

func (r ChatRepository) GetByID(id uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	err := r.DB.First(&chat, "id = ?", id).Error
	if err != nil {
		log.Printf("Failed to run GetByID error: %v", err)
		return nil, err
	}
	return &chat, nil
}
