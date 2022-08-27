package controllers

import (
	"encoding/json"
	"fmt"
	"msn/pkg/models"
	"msn/storage"
	"msn/storage/sqlite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestChatController(t *testing.T) {
	db := sqlite.SetupDatabase()
	chatRepository := storage.NewChatRepository(db)
	userRepository := storage.NewUserRepository(db)
	chatController := NewChatController(chatRepository)

	t.Run("case=must-list-chats", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM chats")
		sqlite.DB.Exec("DELETE FROM users")
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
		sqlite.DB.Exec("DELETE FROM chats")
		sqlite.DB.Exec("DELETE FROM users")
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
		sqlite.DB.Exec("DELETE FROM chats")
		sqlite.DB.Exec("DELETE FROM users")
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

	// t.Run("case=must-create-chat", func(t *testing.T) {
	// 	sqlite.DB.Exec("DELETE FROM chats")
	// 	sqlite.DB.Exec("DELETE FROM users")
	// 	var b bytes.Buffer
	// 	chat := models.Chat{}
	// 	err := json.NewEncoder(&b).Encode(chat)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	req := httptest.NewRequest(http.MethodPost, "/chats", &b)
	// 	w := httptest.NewRecorder()

	// 	chatController.CreateChat(w, req)
	// })
}
