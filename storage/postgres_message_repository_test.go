package storage

import (
	"msn/pkg/models"
	"msn/storage/postgres"
	"testing"

	"github.com/google/uuid"
)

func TestPostgresMessageRepository_SaveAttachment(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	messageRepository := NewMessageRepository(db)
	userRepository := NewUserRepository(db)
	cleaner := NewCleaner(db)

	t.Run("case=must-save-new-chat-with-users", func(t *testing.T) {
		cleaner.Clean()

		user := models.User{
			Username: "xxx",
			Password: "xxx",
		}

		err := userRepository.Create(&user)
		if err != nil {
			t.Errorf("Error creating user: %v", err)
		}

		message := models.Message{
			Content:  "xxx",
			ChatID:   uuid.New(),
			SenderID: user.ID,
		}

		err = messageRepository.Create(&message)
		if err != nil {
			t.Errorf("Error creating message: %v", err)
		}

		attachment := models.Attachment{
			Path:      "path",
			MessageID: message.ID,
		}

		err = messageRepository.SaveAttachment(&attachment)

		if err != nil {
			t.Errorf("Error saving attachment: %v", err)
		}
	})

}
