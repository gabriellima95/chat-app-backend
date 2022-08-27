package data

import (
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"msn/pkg/models"
	"msn/storage"
	"msn/storage/sqlite"
)

func Populate() {

	db := sqlite.SetupDatabase()
	userRepository := storage.NewUserRepository(db)
	chatRepository := storage.NewChatRepository(db)
	messageRepository := storage.NewMessageRepository(db)

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
	userRepository.Create(&paulo)
	userRepository.Create(&gabriel)
	userRepository.Create(&matheus)

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
	chatRepository.Create(chatpg)
	chatRepository.Create(chatgm)
	chatRepository.Create(chatmp)

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
	messageRepository.Create(&m1)
	messageRepository.Create(&m2)
	messageRepository.Create(&m3)
	messageRepository.Create(&m4)
	messageRepository.Create(&m5)
	messageRepository.Create(&m6)
}

func Clear() {
	db := sqlite.SetupDatabase()
	db.Exec("DELETE FROM chats")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM messages")
}

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
