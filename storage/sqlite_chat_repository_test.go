package storage

import (
	"msn/pkg/models"
	"msn/storage/sqlite"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSQLiteChatRepository(t *testing.T) {
	db := sqlite.SetupDatabase()
	chatRepository := NewChatRepository(db)
	userRepository := NewUserRepository(db)
	cleaner := Cleaner{db}

	t.Run("case=must-save-new-chat-with-users", func(t *testing.T) {
		cleaner.Clean()
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
		chat := &models.Chat{
			User1:         user1,
			User2:         user2,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
	})

	t.Run("case=must-populate-users-ids-when-passing-only-struct", func(t *testing.T) {
		cleaner.Clean()
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
		chat := &models.Chat{
			User1:         user1,
			User2:         user2,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}

		err := chatRepository.Create(chat)

		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
		if chat.User1ID != user1.ID {
			t.Errorf("Should populate user id")
		}
		if chat.User2ID != user2.ID {
			t.Errorf("Should populate user id")
		}
	})

	// t.Run("case=returns-error-when-passing-unknown-users-id", func(t *testing.T) {
	// 	sqlite.DB.Exec("DELETE FROM chats")
	// 	sqlite.DB.Exec("DELETE FROM users")
	// 	chat := &models.Chat{
	// 		ID:            uuid.New(),
	// 		User1ID:       uuid.New(),
	// 		User2ID:       uuid.New(),
	// 		LastMessageAt: time.Now(),
	// 		LastMessage:   "oie",
	// 	}

	// 	err := chatRepository.Create(chat)

	// 	if err == nil {
	// 		fmt.Println("")
	// 		t.Errorf("Error saving chat: Should not save with non nullable fields as nil")
	// 	}
	// 	fmt.Println("")
	// })

	t.Run("case=must-not-save-chat-with-non-nullable-fields-as-nil", func(t *testing.T) {
		cleaner.Clean()
		chat := &models.Chat{
			// User1ID:       uuid.New(),
			// User2ID:       uuid.New(),
			// LastMessageAt: time.Now(),
			// LastMessage:   "oie",
		}

		err := chatRepository.Create(chat)

		if err == nil {
			t.Errorf("Error saving chat: Should not save with non nullable fields as nil")
		}

	})

	t.Run("case=does-not-populate-users-when-passing-only-ids", func(t *testing.T) {
		cleaner.Clean()
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
		chat := &models.Chat{
			User1ID:       user1.ID,
			User2ID:       user2.ID,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}

		err := chatRepository.Create(chat)

		emptyUser := models.User{}
		if err != nil {
			t.Errorf("Error saving chat: %v", err)
		}
		if chat.User1.ID != emptyUser.ID {
			t.Errorf("Should not populate users")
		}
		if chat.User2.ID != emptyUser.ID {
			t.Errorf("Should not populate users")
		}
	})
	// TESTAR FOREIGN KEY -> PASSAR USERID QUE NÃO TEM USER ATRELADO

	t.Run("case=must-list-chats-with-matching-user-id", func(t *testing.T) {
		cleaner.Clean()
		user := models.User{
			Username: "111",
			Password: "111",
		}
		userRepository.Create(&user)
		chat1 := &models.Chat{
			User1ID:       user.ID,
			User2ID:       uuid.New(),
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chat2 := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       user.ID,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chat3 := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat1)
		chatRepository.Create(chat2)
		chatRepository.Create(chat3)

		chats, err := chatRepository.ListByUserID(user.ID)

		if err != nil {
			t.Errorf("Error listing chats")
		}

		if !containsChat(*chat1, chats) {
			t.Errorf("List should contain %v", chat1)
		}
		if !containsChat(*chat2, chats) {
			t.Errorf("List should contain %v", chat2)
		}
		if containsChat(*chat3, chats) {
			t.Errorf("List should not contain %v", chat3)
		}
		if len(chats) != 2 {
			t.Errorf("Should return list of lenght 2")
		}

	})

	t.Run("case=must-return-empty-list-when-no-chats-are-found-with-matching-user-id", func(t *testing.T) {
		cleaner.Clean()

		chats, err := chatRepository.ListByUserID(uuid.New())

		if err != nil {
			t.Errorf("Error listing chats")
		}
		if len(chats) != 0 {
			t.Errorf("Should return list of lenght 0")
		}

	})

	t.Run("case=must-return-chats-with-user-fields-populated", func(t *testing.T) {
		cleaner.Clean()
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
		chat := &models.Chat{
			User1ID:       user1.ID,
			User2ID:       user2.ID,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)

		chats, err := chatRepository.ListByUserID(user1.ID)

		if err != nil {
			t.Errorf("Error listing chats: %v", err)
		}
		if !usersMatch(chats[0].User1, user1) {
			t.Errorf("Should load User1")
		}
		if !usersMatch(chats[0].User2, user2) {
			t.Errorf("Should load User2")
		}
	})

	t.Run("case=must-return-chats-ordered-by-last-message-at", func(t *testing.T) {
		cleaner.Clean()
		userID := uuid.New()
		timestamp := time.Now()
		beforeChat := &models.Chat{
			User1ID:       userID,
			User2ID:       uuid.New(),
			LastMessageAt: timestamp,
			LastMessage:   "oie",
		}
		afterChat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       userID,
			LastMessageAt: timestamp.Add(1 * time.Second),
			LastMessage:   "oie",
		}
		chatRepository.Create(beforeChat)
		chatRepository.Create(afterChat)

		chats, err := chatRepository.ListByUserID(userID)

		if err != nil {
			t.Errorf("Error listing chats: %v", err)
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

	t.Run("case=must-update-chat", func(t *testing.T) {
		cleaner.Clean()
		userID := uuid.New()
		timestamp := time.Now()
		newTimestamp := time.Now()
		newLastMessage := "tchau"
		chat := &models.Chat{
			User1ID:       userID,
			User2ID:       uuid.New(),
			LastMessageAt: timestamp,
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		chat.LastMessageAt = newTimestamp
		chat.LastMessage = newLastMessage

		err := chatRepository.Update(chat)

		if err != nil {
			t.Errorf("Error updating chat: %v", err)
		}
		savedChats, err := chatRepository.ListByUserID(userID)
		if err != nil {
			t.Errorf("Error listing chats for assertion: %v", err)
		}
		savedChat := savedChats[0]
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

	t.Run("case=must-get-chat-by-id", func(t *testing.T) {
		cleaner.Clean()
		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
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
}

func containsChat(target models.Chat, chats []models.Chat) bool {
	for _, chat := range chats {
		if chat.ID == target.ID {
			return true
		}
	}
	return false
}

func usersMatch(user1, user2 models.User) bool {
	if user1.ID != user2.ID {
		return false
	}
	if user1.Username != user2.Username {
		return false
	}
	if user1.Password != user2.Password {
		return false
	}

	return true
}
