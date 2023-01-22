package storage

import (
	"msn/pkg/models"
	"msn/storage/sqlite"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMessageRepository(t *testing.T) {
	db := sqlite.SetupDatabase()
	messageRepository := NewMessageRepository(db)
	userRepository := NewUserRepository(db)
	cleaner := NewCleaner(db)

	t.Run("case=must-save-new-message", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		message := &models.Message{
			ID:       uuid.New(),
			Content:  "oie",
			ChatID:   uuid.New(),
			SenderID: uuid.New(),
		}
		now := time.Now()

		err := messageRepository.Create(message)

		if err != nil {
			t.Errorf("Error saving message: %v", err)
		}
		if message.CreatedAt.Before(now) {
			t.Errorf("Error saving message: created_at incorrect")
		}
		if message.UpdatedAt.Before(now) {
			t.Errorf("Error saving message: created_at incorrect")
		}
	})

	t.Run("case=must-not-save-message-with-non-nullable-fields-as-nil", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		message := &models.Message{
			ID: uuid.New(),
		}

		err := messageRepository.Create(message)

		if err == nil {
			t.Errorf("Error saving message: Should not save with non nullable fields as nil")
		}
	})

	t.Run("case=must-list-messages-with-matching-chat-id", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		savedMessage := &models.Message{
			ID:       uuid.New(),
			Content:  "oie",
			ChatID:   uuid.New(),
			SenderID: uuid.New(),
		}
		messageRepository.Create(savedMessage)

		messages, err := messageRepository.ListByChatID(savedMessage.ChatID)

		if err != nil {
			t.Errorf("Error listing chats %v", err)
		}

		if len(messages) != 1 {
			t.Errorf("Should list one message %v", messages)
		}
		if messages[0].ID != savedMessage.ID {
			t.Errorf("Should list saved message %v", messages)
		}
	})

	t.Run("case=must-list-messages-returning-the-user", func(t *testing.T) {
		cleaner.Clean()

		user := &models.User{
			Username: "111",
			Password: "111",
		}
		userRepository.Create(user)
		message := &models.Message{
			Content:  "oie",
			ChatID:   uuid.New(),
			SenderID: user.ID,
		}
		messageRepository.Create(message)

		messages, err := messageRepository.ListByChatID(message.ChatID)

		if err != nil {
			t.Errorf("Error listing chats")
		}

		if len(messages) != 1 {
			t.Errorf("Should list one message %v", messages)
		}
		if messages[0].Sender.Username != user.Username {
			t.Errorf("Should list saved message %v with populated user username %s", messages, user.Username)
		}
		if messages[0].Sender.Password != user.Password {
			t.Errorf("Should list saved message %v with populated user password %s", messages, user.Password)
		}
		if messages[0].Sender.ID != user.ID {
			t.Errorf("Should list saved message %v with populated user id %s", messages, user.ID.String())
		}
	})

	t.Run("case=must-return-empty-list-when-no-chats-are-found-with-matching-user-id", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")

		messages, err := messageRepository.ListByChatID(uuid.New())

		if err != nil {
			t.Errorf("Error listing chats")
		}
		if len(messages) != 0 {
			t.Errorf("Should return list of lenght 0")
		}

	})

	t.Run("case=must-list-messages-ordered-by-created-at-desc", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		chatID := uuid.New()
		firstSavedMessage := &models.Message{
			ID:       uuid.New(),
			Content:  "first message",
			ChatID:   chatID,
			SenderID: uuid.New(),
		}
		messageRepository.Create(firstSavedMessage)
		time.Sleep(time.Second * 1)
		secondSavedMessage := &models.Message{
			ID:       uuid.New(),
			Content:  "second message",
			ChatID:   chatID,
			SenderID: uuid.New(),
		}
		messageRepository.Create(secondSavedMessage)

		messages, err := messageRepository.ListByChatID(chatID)

		if err != nil {
			t.Errorf("Error listing chats")
		}

		if len(messages) != 2 {
			t.Errorf("Should list two messages %v", messages)
		}
		if !secondSavedMessage.CreatedAt.After(firstSavedMessage.CreatedAt) {
			t.Errorf("Message saved last should be after first saved message %v", messages)
		}
		if messages[0].ID != secondSavedMessage.ID {
			t.Errorf("First message in list should be the one with bigger created_at %v", messages)
		}
		if messages[1].ID != firstSavedMessage.ID {
			t.Errorf("Last message in list should be the one with smaller created_at %v", messages)
		}
	})
}
