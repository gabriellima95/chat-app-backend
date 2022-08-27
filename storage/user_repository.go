package storage

import (
	"log"
	"msn/pkg/models"
	"msn/storage/sqlite"

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
	userStorageModel := sqlite.NewUser(user)

	err := r.DB.Create(userStorageModel).Error
	if err != nil {
		return err
	}

	user.ID = userStorageModel.ID
	user.Username = userStorageModel.Username
	user.Password = userStorageModel.Password
	user.CreatedAt = userStorageModel.CreatedAt
	user.UpdatedAt = userStorageModel.UpdatedAt
	user.DeletedAt = userStorageModel.DeletedAt

	return nil
}

func (r UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	var userStorageModel sqlite.User
	err := r.DB.First(&userStorageModel, "username = ?", username).Error
	if err != nil {
		log.Printf("Failed to run GetByUsername error: %v", err)
		return nil, err
	}

	user.ID = userStorageModel.ID
	user.Username = userStorageModel.Username
	user.Password = userStorageModel.Password
	user.CreatedAt = userStorageModel.CreatedAt
	user.UpdatedAt = userStorageModel.UpdatedAt
	user.DeletedAt = userStorageModel.DeletedAt

	return &user, nil
}
