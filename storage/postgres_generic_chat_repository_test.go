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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
		postgres.DB.Exec("DELETE FROM users")
		chat := &models.GenericChat{}

		err := chatRepository.Create(chat)

		if err == nil {
			t.Errorf("Error saving chat: Should not save with non nullable fields as nil")
		}

	})

	t.Run("case=must-save-new-chat-with-empty-name-field", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		postgres.DB.Exec("DELETE FROM generic_chats")
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

func TestPostgresGenericChatRepositoryListByUserID(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	chatRepository := NewGenericChatRepository(db)

	t.Run("case=must-list-chats-with-matching-user-id", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM generic_chats")
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
		user3 := models.User{
			ID:       uuid.New(),
			Username: "333",
			Password: "333",
		}
		chat1 := &models.GenericChat{
			Name:          "grupo1",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}
		chat2 := &models.GenericChat{
			Name:          "grupo2",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user3},
		}
		chat3 := &models.GenericChat{
			Name:          "grupo3",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user2, user3},
		}
		chatRepository.Create(chat1)
		chatRepository.Create(chat2)
		chatRepository.Create(chat3)

		chats, err := chatRepository.ListByUserID(user1.ID)

		if err != nil {
			t.Errorf("Error listing chats")
		}

		if !containsGenericChat(*chat1, chats) {
			t.Errorf("List should contain %v", chat1)
		}
		if !containsGenericChat(*chat2, chats) {
			t.Errorf("List should contain %v", chat2)
		}
		if containsGenericChat(*chat3, chats) {
			t.Errorf("List should not contain %v", chat3)
		}
		if len(chats) != 2 {
			t.Errorf("Should return list of lenght 2")
		}
	})

	t.Run("case=must-list-chats-with-loaded-users", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM generic_chats")
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
			Name:          "grupo1",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user1, user2},
		}
		chatRepository.Create(chat)

		chats, err := chatRepository.ListByUserID(user1.ID)

		if err != nil {
			t.Errorf("Error listing chats")
		}
		if len(chats) != 1 {
			t.Errorf("Should return list of lenght 1")
		}
		if len(chats[0].Users) != 2 {
			t.Errorf("Should return chats with loaded users")
		}
		if !containsUser(user1, chats[0].Users) {
			t.Errorf("Chat should contain the users")
		}
		if !containsUser(user2, chats[0].Users) {
			t.Errorf("Chat should contain the users")
		}
	})

	t.Run("case=must-return-empty-list-when-there-are-no-chats-with-user-id", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM generic_chats")
		postgres.DB.Exec("DELETE FROM users")

		chats, err := chatRepository.ListByUserID(uuid.New())

		if err != nil {
			t.Errorf("Error listing chats")
		}
		if len(chats) != 0 {
			t.Errorf("Should return list of lenght 0")
		}
	})

	t.Run("case=must-list-chats-with-ordered-by-last-message-at", func(t *testing.T) {
		postgres.DB.Exec("DELETE FROM user_chats")
		postgres.DB.Exec("DELETE FROM generic_chats")
		postgres.DB.Exec("DELETE FROM users")
		user := models.User{
			ID:       uuid.New(),
			Username: "111",
			Password: "111",
		}
		beforeChat := &models.GenericChat{
			Name:          "beforeChat",
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user},
		}
		afterChat := &models.GenericChat{
			Name:          "afterChat",
			LastMessageAt: time.Now().Add(5 * time.Second),
			LastMessage:   "oie",
			IsGroup:       true,
			Users:         []models.User{user},
		}
		chatRepository.Create(beforeChat)
		chatRepository.Create(afterChat)

		chats, err := chatRepository.ListByUserID(user.ID)

		if err != nil {
			t.Errorf("Error listing chats")
		}
		if !chats[0].LastMessageAt.After(chats[1].LastMessageAt) {
			t.Errorf("Chats should be ordered by last message at")
		}
		if chats[0].ID != afterChat.ID {
			t.Errorf("After Chat should be first on list")
		}
		if chats[1].ID != beforeChat.ID {
			t.Errorf("Before Chat should be last on list")
		}
	})
}

func containsGenericChat(target models.GenericChat, chats []models.GenericChat) bool {
	for _, chat := range chats {
		if chat.ID == target.ID {
			return true
		}
	}
	return false
}
