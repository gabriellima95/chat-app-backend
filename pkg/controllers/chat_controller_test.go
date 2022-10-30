package controllers

import (
	"encoding/json"
	"fmt"
	"msn/pkg/models"
	"msn/storage"
	"msn/storage/postgres"
	"msn/storage/sqlite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestChatController(t *testing.T) {
	db := sqlite.SetupDatabase()
	cleaner := storage.NewCleaner(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)
	userRepository := storage.NewUserRepository(db)
	chatController := NewChatController(chatRepository, genericChatRepository)

	t.Run("case=must-list-chats", func(t *testing.T) {
		cleaner.Clean()
		user1 := models.User{
			Username: "111",
			Password: "111",
		}
		userRepository.Create(&user1)
		user2 := models.User{
			Username: "222",
			Password: "222",
		}
		userRepository.Create(&user2)
		chat := &models.Chat{
			User1ID:       user1.ID,
			User2ID:       user2.ID,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		url := "/" + user1.ID.String() + "/chats"
		vars := map[string]string{
			"user_id": user1.ID.String(),
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		chatController.ListChats(w, req)

		responseBody := []map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error listing chats: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		responseChat := responseBody[0]
		if fmt.Sprint(responseChat["id"]) != chat.ID.String() {
			t.Errorf("ID field should be %v", chat.ID.String())
		}
		if fmt.Sprint(responseChat["last_message"]) != chat.LastMessage {
			t.Errorf("LastMessage field should be %v", chat.LastMessage)
		}
		// if fmt.Sprint(responseChat["last_message_at"]) != chat.LastMessage {
		// 	t.Errorf("LastMessage field should be %v", chat.LastMessage)
		// }
		userResponse := responseChat["user"].(map[string]interface{})
		if fmt.Sprint(userResponse["id"]) != user1.ID.String() {
			t.Errorf("user ID field should be %v", user1.ID.String())
		}
		if fmt.Sprint(userResponse["username"]) != user1.Username {
			t.Errorf("Username ID field should be %v", user1.Username)
		}
		contactResponse := responseChat["contact"].(map[string]interface{})
		if fmt.Sprint(contactResponse["id"]) != user2.ID.String() {
			t.Errorf("contact ID field should be %v", user2.ID.String())
		}
		if fmt.Sprint(contactResponse["username"]) != user2.Username {
			t.Errorf("contact Username field should be %v", user2.Username)
		}
	})

	t.Run("case=must-list-chats-with-query-param-user-as-user-and-user2-as-contact", func(t *testing.T) {
		cleaner.Clean()
		user1 := models.User{
			Username: "111",
			Password: "111",
		}
		userRepository.Create(&user1)
		user2 := models.User{
			Username: "222",
			Password: "222",
		}
		userRepository.Create(&user2)
		chat := &models.Chat{
			User1ID:       user1.ID,
			User2ID:       user2.ID,
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		url := "/" + user1.ID.String() + "/chats"
		vars := map[string]string{
			"user_id": user2.ID.String(),
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		chatController.ListChats(w, req)

		responseBody := []map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error listing chats: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		responseChat := responseBody[0]
		if fmt.Sprint(responseChat["id"]) != chat.ID.String() {
			t.Errorf("ID field should be %v", chat.ID.String())
		}
		if fmt.Sprint(responseChat["last_message"]) != chat.LastMessage {
			t.Errorf("LastMessage field should be %v", chat.LastMessage)
		}
		// if fmt.Sprint(responseChat["last_message_at"]) != chat.LastMessage {
		// 	t.Errorf("LastMessage field should be %v", chat.LastMessage)
		// }
		userResponse := responseChat["user"].(map[string]interface{})
		if fmt.Sprint(userResponse["id"]) != user2.ID.String() {
			t.Errorf("user ID field should be %v", user2.ID.String())
		}
		if fmt.Sprint(userResponse["username"]) != user2.Username {
			t.Errorf("Username ID field should be %v", user2.Username)
		}
		contactResponse := responseChat["contact"].(map[string]interface{})
		if fmt.Sprint(contactResponse["id"]) != user1.ID.String() {
			t.Errorf("contact ID field should be %v", user1.ID.String())
		}
		if fmt.Sprint(contactResponse["username"]) != user1.Username {
			t.Errorf("contact Username field should be %v", user1.Username)
		}
	})

	t.Run("case=must-return-error-listing-chats-when-passing-invalid-uuid", func(t *testing.T) {
		cleaner.Clean()
		userID := "123"
		url := "/" + userID + "/chats"
		vars := map[string]string{
			"user_id": userID,
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		chatController.ListChats(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}
	})

}

func TestChatControllerListGenericChats(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	cleaner := storage.NewCleaner(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)

	chatController := NewChatController(chatRepository, genericChatRepository)

	t.Run("case=must-list-chats", func(t *testing.T) {
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
		chat1 := &models.GenericChat{
			LastMessageAt: time.Now().Add(15 * time.Second),
			IsGroup:       false,
			LastMessage:   "chat1",
			LastSenderID:  user1.ID,
			Users:         []models.User{user1, user2},
		}
		chat2 := &models.GenericChat{
			LastMessageAt: time.Now().Add(10 * time.Second),
			IsGroup:       true,
			Name:          "chat name 2",
			LastMessage:   "chat2",
			LastSenderID:  user1.ID,
			Users:         []models.User{user1, user2},
		}
		chat3 := &models.GenericChat{
			LastMessageAt: time.Now().Add(5 * time.Second),
			IsGroup:       true,
			Name:          "chat name 3",
			LastMessage:   "chat3",
			LastSenderID:  user2.ID,
			Users:         []models.User{user1, user2},
		}
		genericChatRepository.Create(chat1)
		genericChatRepository.Create(chat2)
		genericChatRepository.Create(chat3)
		url := "/" + user1.ID.String() + "/generic_chats"
		vars := map[string]string{
			"user_id": user1.ID.String(),
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		chatController.ListGenericChats(w, req)

		responseBody := []map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error listing chats: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		if len(responseBody) != 3 {
			t.Errorf("Response body should list all the correct chats")
		}
		responseChat1 := responseBody[0]
		if fmt.Sprint(responseChat1["id"]) != chat1.ID.String() {
			t.Errorf("ID field for first chat should be %v", chat1.ID.String())
		}
		if fmt.Sprint(responseChat1["last_message"]) != "chat1" {
			t.Errorf("LastMessage field is %v and should be %v", responseChat1["last_message"], "chat1")
		}
		if fmt.Sprint(responseChat1["name"]) != user2.Username {
			t.Errorf("Name field should be %v", user2.Username)
		}
		responseChat2 := responseBody[1]
		if fmt.Sprint(responseChat2["id"]) != chat2.ID.String() {
			t.Errorf("ID field for first chat should be %v", chat2.ID.String())
		}
		if fmt.Sprint(responseChat2["last_message"]) != "Eu: chat2" {
			t.Errorf("LastMessage field is %v and should be %v", responseChat2["last_message"], "Eu: chat2")
		}
		if fmt.Sprint(responseChat2["name"]) != chat2.Name {
			t.Errorf("Name field should be %v", chat2.Name)
		}
		responseChat3 := responseBody[2]
		if fmt.Sprint(responseChat3["id"]) != chat3.ID.String() {
			t.Errorf("ID field for first chat should be %v", chat3.ID.String())
		}
		if fmt.Sprint(responseChat3["last_message"]) != "222: chat3" {
			t.Errorf("LastMessage field is %v and should be %v", responseChat3["last_message"], "222: chat3")
		}
		if fmt.Sprint(responseChat3["name"]) != chat3.Name {
			t.Errorf("Name field should be %v", chat3.Name)
		}
	})

	t.Run("case=must-return-error-listing-chats-when-passing-invalid-uuid", func(t *testing.T) {
		cleaner.Clean()
		userID := "123"
		url := "/" + userID + "/chats"
		vars := map[string]string{
			"user_id": userID,
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		chatController.ListGenericChats(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}

	})
}
