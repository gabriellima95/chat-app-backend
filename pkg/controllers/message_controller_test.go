package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"msn/pkg/models"
	"msn/storage"
	"msn/storage/postgres"
	"msn/storage/sqlite"
	"msn/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestMessageController(t *testing.T) {
	db := sqlite.SetupDatabase()
	messageRepository := storage.NewMessageRepository(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)
	notifierMock := &websocket.NotifierMock{}
	messageController := NewMessageController(messageRepository, chatRepository, genericChatRepository, notifierMock)

	t.Run("case=must-create-message", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		sqlite.DB.Exec("DELETE FROM chats")
		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now(),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		responseBody := map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error creating user: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		if len(fmt.Sprint(responseBody["id"])) == 0 {
			t.Errorf("ID field should be populated")
		}
		if fmt.Sprint(responseBody["content"]) != content {
			t.Errorf("Content field should be %v", content)
		}
		if fmt.Sprint(responseBody["chat_id"]) != chatID {
			t.Errorf("ChatID field should be %v", chatID)
		}
		if fmt.Sprint(responseBody["sender_id"]) != senderID {
			t.Errorf("SenderID field should be %v", senderID)
		}
	})

	t.Run("case=must-return-error-creating-message-when-passing-invalid-uuid-in-chat-id", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")

		content := "content"
		chatID := "123"
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}
	})

	t.Run("case=must-return-error-creating-message-when-passing-invalid-uuid-in-sender-id", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")

		content := "content"
		chatID := uuid.NewString()
		senderID := "123"
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}
	})

	t.Run("case=must-list-messages", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		message := &models.Message{
			ID:       uuid.New(),
			Content:  "Hello",
			ChatID:   uuid.New(),
			SenderID: uuid.New(),
		}
		messageRepository.Create(message)
		url := "/" + message.ChatID.String() + "/messages"
		vars := map[string]string{
			"chat_id": message.ChatID.String(),
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()

		messageController.ListMessages(w, req)

		responseBody := []map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error listing messages: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		responseMessage := responseBody[0]
		if fmt.Sprint(responseMessage["id"]) != message.ID.String() {
			t.Errorf("ID field should be %v", message.ID.String())
		}
		if fmt.Sprint(responseMessage["content"]) != message.Content {
			t.Errorf("Content field should be %v", message.Content)
		}
		if fmt.Sprint(responseMessage["sender_id"]) != message.SenderID.String() {
			t.Errorf("SenderID field should be %v", message.SenderID.String())
		}
		if fmt.Sprint(responseMessage["chat_id"]) != message.ChatID.String() {
			t.Errorf("ChatID field should be %v", message.ChatID.String())
		}
	})

	t.Run("case=after-creating-message-must-update-chat", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		sqlite.DB.Exec("DELETE FROM chats")
		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now().Add(60 * time.Second),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		savedChat, err := chatRepository.GetByID(chat.ID)
		if err != nil {
			t.Errorf("Error fetching chat")
		}
		responseBody := map[string]interface{}{}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Error deserializing response body")
		}
		if savedChat.LastMessage != content {
			t.Errorf("Chat last message should be updated")
		}
		messageCreatedAt, _ := time.Parse(time.RFC3339Nano, responseBody["created_at"].(string))
		if savedChat.LastMessageAt.Format("2006-01-02 15:04:05") != messageCreatedAt.Format("2006-01-02 15:04:05") {
			t.Errorf("Chat last message at should be updated")
		}
	})

	t.Run("case=must-call-notifier-twice-with-correct-params", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM messages")
		sqlite.DB.Exec("DELETE FROM chats")
		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now().Add(60 * time.Second),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		var notifierMessage1 models.Message
		var notifierMessage2 models.Message
		var notifierUserID1 string
		var notifierUserID2 string
		notifierMock.NotifyMessageFn = func(message models.Message, userID string) error {
			if userID == chat.User1ID.String() {
				notifierUserID1 = userID
				notifierMessage1 = message
			}
			if userID == chat.User2ID.String() {
				notifierUserID2 = userID
				notifierMessage2 = message
			}
			return nil
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if notifierUserID1 != chat.User1ID.String() {
			t.Errorf("Must notify to User1ID")
		}
		if notifierUserID2 != chat.User2ID.String() {
			t.Errorf("Must notify to User2ID")
		}
		responseBody := map[string]interface{}{}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Error deserializing response body")
		}
		if notifierMessage1.Content != content || notifierMessage1.SenderID.String() != senderID || notifierMessage1.ChatID.String() != chatID || notifierMessage1.ID.String() != responseBody["id"] {
			t.Errorf("Must notify with correct message params")
		}
		if notifierMessage2.Content != content || notifierMessage2.SenderID.String() != senderID || notifierMessage2.ChatID.String() != chatID || notifierMessage2.ID.String() != responseBody["id"] {
			t.Errorf("Must notify with correct message params")
		}
	})
}

func TestCreateGenericMessage(t *testing.T) {
	postgres.Testing = true
	db := postgres.SetupDatabase()
	cleaner := storage.NewCleaner(db)
	messageRepository := storage.NewMessageRepository(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)
	notifierMock := &websocket.NotifierMock{}
	messageController := NewMessageController(messageRepository, chatRepository, genericChatRepository, notifierMock)

	t.Run("case=must-create-message", func(t *testing.T) {
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
		user3 := models.User{
			ID:       uuid.New(),
			Username: "333",
			Password: "333",
		}
		chat := &models.GenericChat{
			LastMessageAt: time.Now().Add(15 * time.Second),
			IsGroup:       false,
			LastMessage:   "chat1",
			LastSenderID:  user1.ID,
			Users:         []models.User{user1, user2, user3},
		}
		genericChatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := user2.ID.String()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		responseBody := map[string]interface{}{}
		if w.Code != 200 {
			t.Errorf("Error creating user: response status code is not 200")
		}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Response body is not valid json")
		}
		if len(fmt.Sprint(responseBody["id"])) == 0 {
			t.Errorf("ID field should be populated")
		}
		if fmt.Sprint(responseBody["content"]) != content {
			t.Errorf("Content field should be %v", content)
		}
		if fmt.Sprint(responseBody["chat_id"]) != chatID {
			t.Errorf("ChatID field should be %v", chatID)
		}
		if fmt.Sprint(responseBody["sender_id"]) != senderID {
			t.Errorf("SenderID field should be %v", senderID)
		}
	})

	t.Run("case=must-return-error-creating-message-when-passing-invalid-uuid-in-chat-id", func(t *testing.T) {
		cleaner.Clean()

		content := "content"
		chatID := "123"
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}
	})

	t.Run("case=must-return-error-creating-message-when-passing-invalid-uuid-in-sender-id", func(t *testing.T) {
		cleaner.Clean()

		content := "content"
		chatID := uuid.NewString()
		senderID := "123"
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if w.Code != 400 {
			t.Errorf("Error listing chats: response status code is not 400")
		}
	})

	t.Run("case=after-creating-message-must-update-chat", func(t *testing.T) {
		cleaner.Clean()

		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now().Add(60 * time.Second),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		savedChat, err := chatRepository.GetByID(chat.ID)
		if err != nil {
			t.Errorf("Error fetching chat")
		}
		responseBody := map[string]interface{}{}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Error deserializing response body")
		}
		if savedChat.LastMessage != content {
			t.Errorf("Chat last message should be updated")
		}
		messageCreatedAt, _ := time.Parse(time.RFC3339Nano, responseBody["created_at"].(string))
		if savedChat.LastMessageAt.Format("2006-01-02 15:04:05") != messageCreatedAt.Format("2006-01-02 15:04:05") {
			t.Errorf("Chat last message at should be updated")
		}
	})

	t.Run("case=must-call-notifier-twice-with-correct-params", func(t *testing.T) {
		cleaner.Clean()

		chat := &models.Chat{
			User1ID:       uuid.New(),
			User2ID:       uuid.New(),
			LastMessageAt: time.Now().Add(60 * time.Second),
			LastMessage:   "oie",
		}
		chatRepository.Create(chat)
		content := "content"
		chatID := chat.ID.String()
		senderID := uuid.NewString()
		jsonMap := map[string]string{"content": content, "chat_id": chatID, "sender_id": senderID}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(jsonMap)
		if err != nil {
			t.Fatal(err)
		}
		var notifierMessage1 models.Message
		var notifierMessage2 models.Message
		var notifierUserID1 string
		var notifierUserID2 string
		notifierMock.NotifyMessageFn = func(message models.Message, userID string) error {
			if userID == chat.User1ID.String() {
				notifierUserID1 = userID
				notifierMessage1 = message
			}
			if userID == chat.User2ID.String() {
				notifierUserID2 = userID
				notifierMessage2 = message
			}
			return nil
		}

		req := httptest.NewRequest(http.MethodPost, "/message", &b)
		w := httptest.NewRecorder()

		messageController.CreateMessage(w, req)

		if notifierUserID1 != chat.User1ID.String() {
			t.Errorf("Must notify to User1ID")
		}
		if notifierUserID2 != chat.User2ID.String() {
			t.Errorf("Must notify to User2ID")
		}
		responseBody := map[string]interface{}{}
		if err := json.Unmarshal([]byte(w.Body.String()), &responseBody); err != nil {
			t.Errorf("Error deserializing response body")
		}
		if notifierMessage1.Content != content || notifierMessage1.SenderID.String() != senderID || notifierMessage1.ChatID.String() != chatID || notifierMessage1.ID.String() != responseBody["id"] {
			t.Errorf("Must notify with correct message params")
		}
		if notifierMessage2.Content != content || notifierMessage2.SenderID.String() != senderID || notifierMessage2.ChatID.String() != chatID || notifierMessage2.ID.String() != responseBody["id"] {
			t.Errorf("Must notify with correct message params")
		}
	})
}
