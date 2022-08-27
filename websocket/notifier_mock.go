package websocket

import (
	"msn/pkg/models"
	"net/http"

	"github.com/google/uuid"
)

type NotifierMock struct {
	NotifyMessageFn      func(message models.Message, userID string) error
	NotifyMessageCounter int
}

func (n *NotifierMock) AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {}

func (n *NotifierMock) NotifyMessage(message models.Message, userID string) error {
	n.NotifyMessageCounter++
	if n.NotifyMessageFn != nil {
		return n.NotifyMessageFn(message, userID)
	}
	return nil
}
