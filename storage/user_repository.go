package storage

import (
	"log"
	"msn/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		DB: db,
	}
}

func (r UserRepository) Create(user *models.User) error {
	user.ID = uuid.New()
	return r.DB.Create(user).Error
}

func (r UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "username = ?", username).Error
	if err != nil {
		log.Printf("Failed to run GetByUsername error: %v", err)
		return nil, err
	}

	return &user, nil
}
