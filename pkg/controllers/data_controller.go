package controllers

import (
	"encoding/json"
	"msn/pkg/models"
	"msn/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	userID, _ := uuid.Parse("15e26fbf-2a2f-4e77-80dd-acbb5cfa6e35")
	paulo := models.User{
		ID:       userID,
		Username: "paulo",
		Password: hashAndSalt("paulo"),
	}
	gabriel := models.User{
		ID:       uuid.New(),
		Username: "gabriel",
		Password: hashAndSalt("gabriel"),
	}
	matheus := models.User{
		ID:       uuid.New(),
		Username: "matheus",
		Password: hashAndSalt("matheus"),
	}
	d.userRepository.Create(&paulo)
	d.userRepository.Create(&gabriel)
	d.userRepository.Create(&matheus)

	chatpg := &models.Chat{
		ID:            uuid.New(),
		User1ID:       paulo.ID,
		User2ID:       gabriel.ID,
		User1:         paulo,
		User2:         gabriel,
		LastMessageAt: time.Now(),
		LastMessage:   "tchau",
	}

	chatgm := &models.Chat{
		ID:            uuid.New(),
		User1ID:       gabriel.ID,
		User2ID:       matheus.ID,
		User1:         gabriel,
		User2:         matheus,
		LastMessageAt: time.Now(),
		LastMessage:   "oie",
	}

	chatmp := &models.Chat{
		ID:            uuid.New(),
		User1ID:       matheus.ID,
		User2ID:       paulo.ID,
		User1:         matheus,
		User2:         paulo,
		LastMessageAt: time.Now(),
		LastMessage:   "tetete",
	}
	d.chatRepository.Create(chatpg)
	d.chatRepository.Create(chatgm)
	d.chatRepository.Create(chatmp)

	m1 := models.Message{
		ID:       uuid.New(),
		Content:  "fwdgedfg",
		ChatID:   chatpg.ID,
		SenderID: paulo.ID,
	}
	m2 := models.Message{
		ID:       uuid.New(),
		Content:  "dfgsfhfcvnc",
		ChatID:   chatpg.ID,
		SenderID: gabriel.ID,
	}
	m3 := models.Message{
		ID:       uuid.New(),
		Content:  "sdvdfhdf",
		ChatID:   chatgm.ID,
		SenderID: gabriel.ID,
	}
	m4 := models.Message{
		ID:       uuid.New(),
		Content:  "asdfdfhdfgb",
		ChatID:   chatgm.ID,
		SenderID: matheus.ID,
	}
	m5 := models.Message{
		ID:       uuid.New(),
		Content:  "fgdfghdfghd",
		ChatID:   chatmp.ID,
		SenderID: paulo.ID,
	}
	m6 := models.Message{
		ID:       uuid.New(),
		Content:  "fgdfghdfghd",
		ChatID:   chatmp.ID,
		SenderID: matheus.ID,
	}
	d.messageRepository.Create(&m1)
	d.messageRepository.Create(&m2)
	d.messageRepository.Create(&m3)
	d.messageRepository.Create(&m4)
	d.messageRepository.Create(&m5)
	d.messageRepository.Create(&m6)

	json.NewEncoder(w).Encode("")
}

func (d DataController) Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	DB := d.chatRepository.DB
	DB.Exec("DELETE FROM messages")
	DB.Exec("DELETE FROM chats")
	DB.Exec("DELETE FROM users")
}
