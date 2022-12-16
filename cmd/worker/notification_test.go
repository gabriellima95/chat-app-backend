package worker

import (
	"msn/pkg/models"
	"msn/websocket"
	"testing"

	"github.com/google/uuid"
)

func TestNotificationWorker(t *testing.T) {
	notifierMock := &websocket.NotifierMock{}
	ch := make(chan NotificationTask)
	notificationWorker := Notification{notifierMock, ch}

	t.Run("case=must-call-notifier-correctly", func(t *testing.T) {
		var capturedMessage models.Message
		var capturedUserID string
		notifierMock.NotifyMessageFn = func(message models.Message, userID string) error {
			capturedMessage = message
			capturedUserID = userID
			return nil
		}
		id := uuid.New()
		chatID := uuid.New()
		senderID := uuid.New()
		content := "oie"
		userID := uuid.NewString()
		message := models.Message{
			ID:       id,
			Content:  content,
			ChatID:   chatID,
			SenderID: senderID,
		}
		task := NotificationTask{Message: message, UserID: userID}

		go notificationWorker.Run()
		notificationWorker.SendTask(task)

		if capturedMessage.ID != id || capturedMessage.Content != content || capturedMessage.ChatID != chatID || capturedMessage.SenderID != senderID {
			t.Errorf("Error sending notification task: message with incorrect fields")
		}
		if capturedUserID != userID {
			t.Errorf("Error sending notification task: incorrect userID")
		}
		notificationWorker.Stop()
	})
}
