package data

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"msn/pkg/models"
	"msn/storage"
	"msn/storage/postgres"
	"msn/storage/sqlite"
)

func Populate(database string) {
	var db *gorm.DB
	if database == "postgres" {
		db = postgres.SetupDatabase()
	} else {
		db = sqlite.SetupDatabase()
	}

	userRepository := storage.NewUserRepository(db)
	chatRepository := storage.NewChatRepository(db)
	genericChatRepository := storage.NewGenericChatRepository(db)
	messageRepository := storage.NewMessageRepository(db)

	PopulateDB(userRepository, chatRepository, messageRepository, genericChatRepository)
}

func PopulateDB(userRepository storage.UserRepository, chatRepository storage.ChatRepository, messageRepository storage.MessageRepository, genericChatRepository storage.GenericChatRepository) {
	paulo := models.User{
		Username: "paulo",
		Password: hashAndSalt("paulo"),
	}
	gabriel := models.User{
		Username: "gabriel",
		Password: hashAndSalt("gabriel"),
	}
	matheus := models.User{
		Username: "matheus",
		Password: hashAndSalt("matheus"),
	}
	userRepository.Create(&paulo)
	userRepository.Create(&gabriel)
	userRepository.Create(&matheus)

	chatpg := &models.Chat{
		User1ID:       paulo.ID,
		User2ID:       gabriel.ID,
		User1:         paulo,
		User2:         gabriel,
		LastMessageAt: time.Now(),
		LastMessage:   "tchau",
	}

	chatgm := &models.Chat{
		User1ID:       gabriel.ID,
		User2ID:       matheus.ID,
		User1:         gabriel,
		User2:         matheus,
		LastMessageAt: time.Now(),
		LastMessage:   "oie",
	}

	chatmp := &models.Chat{
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
		Content:  "fwdgedfg",
		ChatID:   chatpg.ID,
		SenderID: paulo.ID,
	}
	m2 := models.Message{
		Content:  "dfgsfhfcvnc",
		ChatID:   chatpg.ID,
		SenderID: gabriel.ID,
	}
	m3 := models.Message{
		Content:  "sdvdfhdf",
		ChatID:   chatgm.ID,
		SenderID: gabriel.ID,
	}
	m4 := models.Message{
		Content:  "asdfdfhdfgb",
		ChatID:   chatgm.ID,
		SenderID: matheus.ID,
	}
	m5 := models.Message{
		Content:  "fgdfghdfghd",
		ChatID:   chatmp.ID,
		SenderID: paulo.ID,
	}
	m6 := models.Message{
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

	newChatgroup := models.GenericChat{
		Name:          "azilados",
		LastMessage:   "ow yeah",
		LastSenderID:  paulo.ID,
		LastMessageAt: time.Now(),
		Users:         []models.User{paulo, gabriel, matheus},
	}
	newChatpg := models.GenericChat{
		LastMessage:   "top",
		LastSenderID:  paulo.ID,
		LastMessageAt: time.Now(),
		IsGroup:       false,
		Users:         []models.User{paulo, gabriel},
	}
	newChatpm := models.GenericChat{
		LastMessage:   "massa",
		LastSenderID:  matheus.ID,
		LastMessageAt: time.Now(),
		IsGroup:       false,
		Users:         []models.User{paulo, matheus},
	}
	newChatgm := models.GenericChat{
		LastMessage:   "irado",
		LastSenderID:  gabriel.ID,
		LastMessageAt: time.Now(),
		IsGroup:       false,
		Users:         []models.User{gabriel, matheus},
	}

	genericChatRepository.Create(&newChatgroup)
	genericChatRepository.Create(&newChatpg)
	genericChatRepository.Create(&newChatpm)
	genericChatRepository.Create(&newChatgm)

	mm1 := models.Message{
		Content:  "fwdgedfg",
		ChatID:   newChatpg.ID,
		SenderID: paulo.ID,
	}
	mm2 := models.Message{
		Content:  "dfgsfhfcvnc",
		ChatID:   newChatpg.ID,
		SenderID: gabriel.ID,
	}
	mm3 := models.Message{
		Content:  "sdvdfhdf",
		ChatID:   newChatgm.ID,
		SenderID: gabriel.ID,
	}
	mm4 := models.Message{
		Content:  "asdfdfhdfgb",
		ChatID:   newChatgm.ID,
		SenderID: matheus.ID,
	}
	mm5 := models.Message{
		Content:  "fgdfghdfghd",
		ChatID:   newChatpm.ID,
		SenderID: paulo.ID,
	}
	mm6 := models.Message{
		Content:  "fgdfghdfghd",
		ChatID:   newChatpm.ID,
		SenderID: matheus.ID,
	}
	mm7 := models.Message{
		Content:  "wrtefrgefh",
		ChatID:   newChatgroup.ID,
		SenderID: matheus.ID,
	}
	mm8 := models.Message{
		Content:  "dsfgsgdfs",
		ChatID:   newChatgroup.ID,
		SenderID: gabriel.ID,
	}
	mm9 := models.Message{
		Content:  "gedfbdef",
		ChatID:   newChatgroup.ID,
		SenderID: paulo.ID,
	}
	messageRepository.Create(&mm1)
	messageRepository.Create(&mm2)
	messageRepository.Create(&mm3)
	messageRepository.Create(&mm4)
	messageRepository.Create(&mm5)
	messageRepository.Create(&mm6)
	messageRepository.Create(&mm7)
	messageRepository.Create(&mm8)
	messageRepository.Create(&mm9)
}

func Clear(database string) {
	var db *gorm.DB
	if database == "postgres" {
		db = postgres.SetupDatabase()
	} else {
		db = sqlite.SetupDatabase()
	}
	ClearDB(db)
}

func ClearDB(db *gorm.DB) {
	db.Exec("DELETE FROM chats")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM user_chats")
	db.Exec("DELETE FROM generic_chats")
}

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
