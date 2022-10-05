package storage

import (
	"msn/pkg/models"
	"msn/storage/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPostgresGenericChatRepositoryCreate(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	chatRepository := NewGenericChatRepository(db)
	userRepository := NewUserRepository(db)

	t.Run("case=must-save-new-chat-with-users", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user1 := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		user2 := models.User{
			ID:       uuid.New(),
			Username: "222",
			Password: "222",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
	})

	t.Run("case=must-not-save-chat-with-non-nullable-fields-as-nil", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		chat := &models.GenericChat{}

		err := chatRepository.Create(chat)

		if err == nil {
			t.Errorf("Error saving chat: Should not save with non nullable fields as nil")
		}

	})

	t.Run("case=must-save-new-chat-with-empty-name-field", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		chat := &models.GenericChat{
			LastMessageAt: time.Now(),
			LastMessage:   "chat com nome vazio",
			IsGroup:       true,
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
	})

	t.Run("case=must-save-users-associated-with-chat", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user1 := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		user2 := models.User{
			ID:       uuid.New(),
			Username: "222",
			Password: "222",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
	})

	t.Run("case=must-save-chat-with-existing-users", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user1 := models.User{
			Username: "111",
			Password: "111",
		}

		user2 := models.User{
			Username: "222",
			Password: "222",
		}
		userRepository.Create(&user1)
		userRepository.Create(&user2)
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
	})

	t.Run("case=must-save-new-users-along-with-chat", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user},
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
		savedUser, err := userRepository.GetByUsername(user.Username)
		if err != nil {
			t.Errorf("Error fetching saved user: %v", err)
		}
		if savedUser.ID != user.ID {
			t.Errorf("Error users IDs should be the same")
		}
	})
}

func TestPostgresGenericChatRepositoryGetByID(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	chatRepository := NewGenericChatRepository(db)
	// userRepository := NewUserRepository(db)

	t.Run("case=must-get-chat-by-id", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user1 := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		user2 := models.User{
			ID:       uuid.New(),
			Username: "222",
			Password: "222",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}
		chatRepository.Create(chat)

		savedChat, err := chatRepository.GetByID(chat.ID)

		if err != nil {
			t.Errorf("Error getting chat by id: %v", err)
		}
		if savedChat.ID != chat.ID {
			t.Errorf("Chat ID should be the same of query")
		}
	})

	t.Run("case=must-get-chat-by-id-with-users-loaded", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")
		user1 := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		user2 := models.User{
			ID:       uuid.New(),
			Username: "222",
			Password: "222",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}
		chatRepository.Create(chat)

		savedChat, err := chatRepository.GetByID(chat.ID)

		if err != nil {
			t.Errorf("Error getting chat by id: %v", err)
		}
		if savedChat.ID != chat.ID {
			t.Errorf("Chat ID should be the same of query")
		}
		if len(savedChat.Users) != 2 {
			t.Errorf("Chat should be fetched with users loaded")
		}
		if !containsUser(user1, savedChat.Users) {
			t.Errorf("Chat should contains the users")
		}
		if !containsUser(user2, savedChat.Users) {
			t.Errorf("Chat should contains the users")
		}
	})
}

func containsUser(target models.User, users []models.User) bool {
	for _, user := range users {
		if user.ID == target.ID {
			return true
		}
	}
	return false
}

func TestPostgresGenericChatRepositoryUpdate(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	chatRepository := NewGenericChatRepository(db)
	// userRepository := NewUserRepository(db)

	t.Run("case=must-update-chat", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM refactor_chats")
		postgres.DB.Exec("DELETE FROM users")

		timestamp := time.Now()
		newTimestamp := time.Now()
		newLastMessage := "tchau"
		user := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		chat := &models.GenericChat{
			Name:          "grupo",
			LastMessageAt: timestamp,
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user},
		}
		chatRepository.Create(chat)
		chat.LastMessageAt = newTimestamp
		chat.LastMessage = newLastMessage

		err := chatRepository.Update(chat)

		if err != nil {
			t.Errorf("Error updating chat: %v", err)
		}
		savedChat, err := chatRepository.GetByID(chat.ID)
		if err != nil {
			t.Errorf("Error fetching chat for assertion: %v", err)
		}
		if savedChat.ID != chat.ID {
			t.Errorf("Chat must be the same for update")
		}
		if savedChat.LastMessage != newLastMessage {
			t.Errorf("Chat last message should be updated")
		}
		if savedChat.LastMessageAt.Format("2006-01-02 15:04:05") != newTimestamp.Format("2006-01-02 15:04:05") {
			t.Errorf("Chat last message at should be updated")
		}
	})
}
