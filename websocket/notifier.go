package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"msn/pkg/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Notifier interface {
	NotifyMessage(message models.Message, userID string) error
	AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID)
}

type SocketNotifier struct {
	upgrader *websocket.Upgrader
	clients  map[string]*websocket.Conn
}

func NewSocketNotifier() *SocketNotifier {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	clients := make(map[string]*websocket.Conn, 0)

	return &SocketNotifier{
		upgrader: upgrader,
		clients:  clients,
	}
}

func (s *SocketNotifier) AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	userIDstr := userID.String()
	s.clients[userIDstr] = conn
}

func (s *SocketNotifier) NotifyMessage(message models.Message, userID string) error {
	fmt.Println("NotifyMessage", s.clients)

	conn, ok := s.clients[userID]
	if !ok {
		return fmt.Errorf("user-socket-not-connected")
	}

	messageMap := make(map[string]interface{})
	messageMap["id"] = message.ID.String()
	messageMap["chat_id"] = message.ChatID.String()
	messageMap["sender_id"] = message.SenderID.String()
	messageMap["created_at"] = message.CreatedAt
	messageMap["content"] = message.Content

	jsonStr, _ := json.Marshal(messageMap)

	err := conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Chegou e deu bom")
	return nil
}
