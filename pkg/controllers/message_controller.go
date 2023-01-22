package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"msn/pkg/models"
	"msn/storage"
	"msn/websocket"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MessageController struct {
	messageRepository     storage.MessageRepository
	socketNotifier        websocket.Notifier
	chatRepository        storage.ChatRepository
	genericChatRepository storage.GenericChatRepository
}

func NewMessageController(
	messageRepository storage.MessageRepository,
	chatRepository storage.ChatRepository,
	genericChatRepository storage.GenericChatRepository,
	socketNotifier websocket.Notifier,
) MessageController {
	return MessageController{
		messageRepository:     messageRepository,
		chatRepository:        chatRepository,
		socketNotifier:        socketNotifier,
		genericChatRepository: genericChatRepository,
	}
}

func (m MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messagePayload MessageRequestSchema
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	err := json.NewDecoder(r.Body).Decode(&messagePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(messagePayload)

	chatID, err := uuid.Parse(messagePayload.ChatID)
	if err != nil {
		http.Error(w, "chatID should be valid uuid", http.StatusBadRequest)
		return
	}
	senderID, err := uuid.Parse(messagePayload.SenderID)
	if err != nil {
		http.Error(w, "senderID should be valid uuid", http.StatusBadRequest)
		return
	}

	message := models.Message{
		// ID:       uuid.New(),
		Content:  messagePayload.Content,
		ChatID:   chatID,
		SenderID: senderID,
	}

	err = m.messageRepository.Create(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chat, err := m.chatRepository.GetByID(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chat.LastMessage = message.Content
	chat.LastMessageAt = message.CreatedAt

	err = m.chatRepository.Update(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageResponse := MessageResponseSchema{
		ID:        message.ID.String(),
		CreatedAt: message.CreatedAt,
		Content:   message.Content,
		ChatID:    message.ChatID.String(),
		SenderID:  message.SenderID.String(),
	}

	m.socketNotifier.NotifyMessage(message, chat.User1ID.String())
	m.socketNotifier.NotifyMessage(message, chat.User2ID.String())

	json.NewEncoder(w).Encode(messageResponse)
}

func (m MessageController) CreateGenericMessage(w http.ResponseWriter, r *http.Request) {
	var messagePayload MessageRequestSchema
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	err := json.NewDecoder(r.Body).Decode(&messagePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(messagePayload)

	chatID, err := uuid.Parse(messagePayload.ChatID)
	if err != nil {
		http.Error(w, "chatID should be valid uuid", http.StatusBadRequest)
		return
	}
	senderID, err := uuid.Parse(messagePayload.SenderID)
	if err != nil {
		http.Error(w, "senderID should be valid uuid", http.StatusBadRequest)
		return
	}

	message := models.Message{
		// ID:       uuid.New(),
		Content:  messagePayload.Content,
		ChatID:   chatID,
		SenderID: senderID,
	}

	err = m.messageRepository.Create(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chat, err := m.genericChatRepository.GetByID(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chat.LastMessage = message.Content
	chat.LastSenderID = senderID
	chat.LastMessageAt = message.CreatedAt

	err = m.genericChatRepository.Update(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messageResponse := MessageResponseSchema{
		ID:        message.ID.String(),
		CreatedAt: message.CreatedAt,
		Content:   message.Content,
		ChatID:    message.ChatID.String(),
		SenderID:  message.SenderID.String(),
	}

	for _, user := range chat.Users {
		err = m.socketNotifier.NotifyMessage(message, user.ID.String())
		log.Printf("Failed on NotifyMessage for user %s: %v", user.ID.String(), err)
	}

	json.NewEncoder(w).Encode(messageResponse)
}

func (m MessageController) ListMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	ID := params["chat_id"]

	chatID, err := uuid.Parse(ID)
	if err != nil {
		http.Error(w, "chatID should be valid uuid", http.StatusBadRequest)
		return
	}

	messages, err := m.messageRepository.ListByChatID(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messageListResponse []MessageResponseSchema
	for _, message := range messages {
		m := MessageResponseSchema{
			ID:         message.ID.String(),
			CreatedAt:  message.CreatedAt,
			Content:    message.Content,
			ChatID:     message.ChatID.String(),
			SenderID:   message.SenderID.String(),
			SenderName: message.Sender.Username,
		}
		messageListResponse = append(messageListResponse, m)
	}

	json.NewEncoder(w).Encode(messageListResponse)
}
