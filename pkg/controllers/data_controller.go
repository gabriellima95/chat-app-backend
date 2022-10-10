package controllers

import (
	"encoding/json"
	"msn/cmd/data"
	"msn/storage"
	"net/http"
)

type DataController struct {
	chatRepository    storage.ChatRepository
	userRepository    storage.UserRepository
	messageRepository storage.MessageRepository
}

func NewDataController(chatRepository storage.ChatRepository, userRepository storage.UserRepository, messageRepository storage.MessageRepository) DataController {
	return DataController{
		chatRepository:    chatRepository,
		userRepository:    userRepository,
		messageRepository: messageRepository,
	}
}

func (d DataController) Populate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data.PopulateDB(d.userRepository, d.chatRepository, d.messageRepository)

	json.NewEncoder(w).Encode("")
}

func (d DataController) Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	DB := d.chatRepository.DB
	data.ClearDB(DB)
}
