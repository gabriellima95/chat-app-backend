package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageNotification struct {
	ID          string
	ChatID      string
	SenderID    string
	CreatedAt   time.Time
	Content     string
	ChatContent string
	SenderName  string
}

// messageMap["id"] = message.ID.String()
// messageMap["chat_id"] = message.ChatID.String()
// messageMap["sender_id"] = message.SenderID.String()
// messageMap["created_at"] = message.CreatedAt
// messageMap["content"] = message.Content
// messageMap["sender_name"] = message.Sender.Username

type Notifier interface {
	NotifyMessage(message MessageNotification, userID string) error
	AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID)
}

type SocketNotifier struct {
	upgrader *websocket.Upgrader
	connPool map[string]*websocket.Conn
}

func NewSocketNotifier() *SocketNotifier {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	connPool := make(map[string]*websocket.Conn, 0)

	return &SocketNotifier{
		upgrader: upgrader,
		connPool: connPool,
	}
}

func (s *SocketNotifier) AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	userIDstr := userID.String()
	s.connPool[userIDstr] = conn
}

func (s *SocketNotifier) NotifyMessage(message MessageNotification, userID string) error {
	fmt.Println("NotifyMessage", s.connPool)

	conn, ok := s.connPool[userID]
	if !ok {
		return fmt.Errorf("user-socket-not-connected")
	}

	messageMap := make(map[string]interface{})
	messageMap["id"] = message.ID
	messageMap["chat_id"] = message.ChatID
	messageMap["sender_id"] = message.SenderID
	messageMap["created_at"] = message.CreatedAt
	messageMap["content"] = message.Content
	messageMap["chat_content"] = message.ChatContent
	messageMap["sender_name"] = message.SenderName

	jsonStr, _ := json.Marshal(messageMap)

	err := conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Chegou e deu bom")
	return nil
}
