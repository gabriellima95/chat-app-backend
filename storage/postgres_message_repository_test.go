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

func TestPostgresMessageRepository_ListByChatID(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	messageRepository := NewMessageRepository(db)
	userRepository := NewUserRepository(db)
	cleaner := NewCleaner(db)

	t.Run("case=must-list-messages-by-chat-with-attachments", func(t *testing.T) {
		cleaner.Clean()

		user := models.User{
			Username: "xxx",
			Password: "xxx",
		}

		err := userRepository.Create(&user)
		if err != nil {
			t.Errorf("Error creating user: %v", err)
		}

		chatID := uuid.New()
		message1 := models.Message{
			Content:  "xxx",
			ChatID:   chatID,
			SenderID: user.ID,
		}

		message2 := message1

		err = messageRepository.Create(&message1)
		if err != nil {
			t.Errorf("Error creating message: %v", err)
		}

		err = messageRepository.Create(&message2)
		if err != nil {
			t.Errorf("Error creating message: %v", err)
		}

		attachment1 := models.Attachment{
			Path:      "path1",
			MessageID: message1.ID,
		}
		attachment2 := models.Attachment{
			Path:      "path2",
			MessageID: message2.ID,
		}
		attachment3 := models.Attachment{
			Path:      "path3",
			MessageID: message2.ID,
		}

		err = messageRepository.SaveAttachment(&attachment1)
		if err != nil {
			t.Errorf("Error saving attachment: %v", err)
		}
		err = messageRepository.SaveAttachment(&attachment2)
		if err != nil {
			t.Errorf("Error saving attachment: %v", err)
		}
		err = messageRepository.SaveAttachment(&attachment3)
		if err != nil {
			t.Errorf("Error saving attachment: %v", err)
		}

		messageList, err := messageRepository.ListByChatID(chatID)
		if err != nil {
			t.Errorf("Error listing messages: %v", err)
		}
		if len(messageList) != 2 {
			t.Errorf("Error listing messages: list length should be 2")
		}
		for _, message := range messageList {
			if message.ID == message1.ID {
				if len(message.Attachments) != 1 {
					t.Errorf("Error listing messages: attachments length should be 1 for message1")
				}

				if message.Attachments[0].Path != "path1" {
					t.Errorf("Error listing messages: attachment path should be path1 for message1")
				}
			}
			if message.ID == message2.ID {
				if len(message.Attachments) != 2 {
					t.Errorf("Error listing messages: attachments length should be 2 for message2")
				}
			}
		}
	})

}
