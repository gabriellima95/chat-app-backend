package subscribers

import (
	"encoding/json"
	"log"
	"msn/pubsub"
	"msn/storage"
	"msn/websocket"
	"time"

	"github.com/google/uuid"
)

type SendMessageNotificationSubscriber struct {
	socketNotifier websocket.Notifier
	chatRepository storage.ChatRepository
}

func NewSendMessageNotificationSubscriber(socketNotifier websocket.Notifier, chatRepository storage.ChatRepository) SendMessageNotificationSubscriber {
	return SendMessageNotificationSubscriber{
		socketNotifier: socketNotifier,
		chatRepository: chatRepository,
	}
}

type SendMessageNotificationRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
}

func (s SendMessageNotificationSubscriber) Name() string {
	return "send-message-notification-subscriber"
}

func (s SendMessageNotificationSubscriber) Run(msg pubsub.Message) error {
	payloadBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		log.Printf("Unable to Marshal message payload: %s", err.Error())
		return err
	}

	var request SendMessageNotificationRequest
	err = json.Unmarshal(payloadBytes, &request)
	if err != nil {
		log.Printf("Unable to Unmarshal message payload: %s", err.Error())
		return err
	}

	chatUUID, err := uuid.Parse(request.ChatID)
	if err != nil {
		log.Printf("Unable to parse ChatID into UUID: %s", err.Error())
		return err
	}

	chat, err := s.chatRepository.GetByID(chatUUID)
	if err != nil {
		log.Printf("Error on chatRepository.GetByID: %s", err.Error())
		return err
	}

	messageNotification := websocket.MessageNotification{
		ID:          request.ID,
		ChatID:      chat.ID.String(),
		SenderID:    request.SenderID,
		CreatedAt:   request.CreatedAt,
		Content:     request.Content,
		ChatContent: request.Content,
	}
	log.Printf("Message Notification: %s", messageNotification)

	s.socketNotifier.NotifyMessage(messageNotification, chat.User1ID.String())
	s.socketNotifier.NotifyMessage(messageNotification, chat.User2ID.String())
	return nil
}
