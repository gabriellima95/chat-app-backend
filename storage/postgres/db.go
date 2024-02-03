package postgres

import (
	"fmt"
	"msn/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Testing bool = false

func SetupDatabase() *gorm.DB {
	host := "db"
	if Testing {
		host = "localhost"
	}
	user := "root"
	password := "root"
	dbname := "chat"
	if Testing {
		dbname = "chat_test"
	}
	port := "5432"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Chat{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Attachment{})
	db.AutoMigrate(&models.GenericChat{})

	// db.AutoMigrate(&models.Message{})

	DB = db
	return db
}
