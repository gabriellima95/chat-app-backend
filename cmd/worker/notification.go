package worker

import (
	"msn/pkg/models"
	"msn/websocket"
)

type NotificationTask struct {
	Message models.Message
	UserID  string
}

type Notification struct {
	socketNotifier websocket.Notifier
	ch             chan NotificationTask
}

func NewNotificationWorker(socketNotifier websocket.Notifier) Notification {
	ch := make(chan NotificationTask)
	return Notification{socketNotifier, ch}
}

func (w Notification) Run() {
	for task := range w.ch {
		w.socketNotifier.NotifyMessage(task.Message, task.UserID)
	}
}

func (w Notification) SendTask(task NotificationTask) {
	w.ch <- task
}

func (w Notification) Stop() {
	close(w.ch)
}
