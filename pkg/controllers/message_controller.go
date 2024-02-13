package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"msn/pkg/models"
	"msn/pubsub"
	"msn/storage"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MessageController struct {
	messageRepository     storage.MessageRepository
	chatRepository        storage.ChatRepository
	genericChatRepository storage.GenericChatRepository
	storageCLient         storage.FileStorageClient
	publisher             pubsub.Publisher
}

func NewMessageController(
	messageRepository storage.MessageRepository,
	chatRepository storage.ChatRepository,
	genericChatRepository storage.GenericChatRepository,
	publisher pubsub.Publisher,
	storageCLient storage.FileStorageClient,
) MessageController {
	return MessageController{
		messageRepository:     messageRepository,
		chatRepository:        chatRepository,
		publisher:             publisher,
		genericChatRepository: genericChatRepository,
		storageCLient:         storageCLient,
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

	m.publishMessageCreatedEvent(message, messagePayload)

	messageResponse := MessageResponseSchema{
		ID:        message.ID.String(),
		CreatedAt: message.CreatedAt,
		Content:   message.Content,
		ChatID:    message.ChatID.String(),
		SenderID:  message.SenderID.String(),
	}

	json.NewEncoder(w).Encode(messageResponse)
}

func (m MessageController) publishMessageCreatedEvent(message models.Message, request MessageRequestSchema) {
	pubsubMessagePayload := make(map[string]interface{}, 0)
	pubsubMessagePayload["id"] = message.ID
	pubsubMessagePayload["created_at"] = message.CreatedAt
	pubsubMessagePayload["content"] = message.Content
	pubsubMessagePayload["chat_id"] = message.ChatID
	pubsubMessagePayload["sender_id"] = message.SenderID

	var pubsubMessageAttachmentsPayload []map[string]interface{}
	for _, attachment := range request.Attachments {
		pubsubMessageAttachmentsPayload = append(pubsubMessageAttachmentsPayload, map[string]interface{}{
			"content_type": attachment.ContentType,
			"content":      attachment.Content,
			"filename":     attachment.Filename,
		})
	}

	pubsubMessagePayload["attachments"] = pubsubMessageAttachmentsPayload

	pubsubMessage := pubsub.Message{
		Topic:   "message-created",
		Payload: pubsubMessagePayload,
	}

	err := m.publisher.Publish(pubsubMessage)
	if err != nil {
		log.Printf("error on MessageController.publisher.Publish: %s", err.Error())
	}
}

// func (m MessageController) CreateGenericMessage(w http.ResponseWriter, r *http.Request) {
// 	var messagePayload MessageRequestSchema
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "*")

// 	err := json.NewDecoder(r.Body).Decode(&messagePayload)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println(messagePayload)

// 	chatID, err := uuid.Parse(messagePayload.ChatID)
// 	if err != nil {
// 		http.Error(w, "chatID should be valid uuid", http.StatusBadRequest)
// 		return
// 	}
// 	senderID, err := uuid.Parse(messagePayload.SenderID)
// 	if err != nil {
// 		http.Error(w, "senderID should be valid uuid", http.StatusBadRequest)
// 		return
// 	}

// 	message := models.Message{
// 		// ID:       uuid.New(),
// 		Content:  messagePayload.Content,
// 		ChatID:   chatID,
// 		SenderID: senderID,
// 	}

// 	err = m.messageRepository.Create(&message)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	chat, err := m.genericChatRepository.GetByID(chatID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	for _, user := range chat.Users {
// 		if message.SenderID == user.ID {
// 			message.Sender = user
// 		}
// 	}

// 	chat.LastMessage = message.Content
// 	chat.LastSenderID = senderID
// 	chat.LastMessageAt = message.CreatedAt

// 	err = m.genericChatRepository.Update(chat)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	messageResponse := MessageResponseSchema{
// 		ID:        message.ID.String(),
// 		CreatedAt: message.CreatedAt,
// 		Content:   message.Content,
// 		ChatID:    message.ChatID.String(),
// 		SenderID:  message.SenderID.String(),
// 	}

// 	for _, user := range chat.Users {
// 		messageNotification := websocket.MessageNotification{
// 			ID:          message.ID.String(),
// 			ChatID:      chat.ID.String(),
// 			SenderID:    senderID.String(),
// 			CreatedAt:   message.CreatedAt,
// 			Content:     message.Content,
// 			ChatContent: chat.GetLastMessage(user.ID),
// 			SenderName:  message.Sender.Username,
// 		}
// 		err = m.socketNotifier.NotifyMessage(messageNotification, user.ID.String())
// 		log.Printf("Failed on NotifyMessage for user %s: %v", user.ID.String(), err)
// 	}

// 	json.NewEncoder(w).Encode(messageResponse)
// }

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
			ID:        message.ID.String(),
			CreatedAt: message.CreatedAt,
			Content:   message.Content,
			ChatID:    message.ChatID.String(),
			SenderID:  message.SenderID.String(),
			Attachments: []AttachmentResponse{
				{Url: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/132.png", Filename: "image.png"},
			},
		}
		messageListResponse = append(messageListResponse, m)
	}

	json.NewEncoder(w).Encode(messageListResponse)
}
