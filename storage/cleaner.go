package storage

import "gorm.io/gorm"

type Cleaner struct {
	DB *gorm.DB
}

func NewCleaner(db *gorm.DB) Cleaner {
	return Cleaner{
		DB: db,
	}
}

func (c Cleaner) Clean() {
	c.DB.Exec("DELETE FROM chats")
	c.DB.Exec("DELETE FROM user_chats")
	c.DB.Exec("DELETE FROM generic_chats")
	c.DB.Exec("DELETE FROM messages")
	c.DB.Exec("DELETE FROM users")
}
