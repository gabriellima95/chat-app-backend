package websocket

import (
	"net/http"

	"github.com/google/uuid"
)

type NotifierMock struct {
	NotifyMessageFn      func(message MessageNotification, userID string) error
	NotifyMessageCounter int
}

func (n *NotifierMock) AddConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {}

func (n *NotifierMock) NotifyMessage(message MessageNotification, userID string) error {
	n.NotifyMessageCounter++
	if n.NotifyMessageFn != nil {
		return n.NotifyMessageFn(message, userID)
	}
	return nil
}

func (n *NotifierMock) NotifyAttachment(message AttachmentNotification, userID string) error {
	return nil
}
