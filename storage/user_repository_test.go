package storage

import (
	"msn/pkg/models"
	"msn/storage/sqlite"
	"testing"
	"time"
)

func TestUserRepository(t *testing.T) {
	db := sqlite.SetupDatabase()
	userRepository := NewUserRepository(db)

	t.Run("case=must-save-new-user", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		user := &models.User{
			Username: "abc",
			Password: "abc",
		}
		now := time.Now()

		err := userRepository.Create(user)

		if err != nil {
			t.Errorf("Error saving user: %v", err)
		}
		if user.CreatedAt.Before(now) {
			t.Errorf("Error saving user: created_at incorrect")
		}
		if user.UpdatedAt.Before(now) {
			t.Errorf("Error saving user: created_at incorrect")
		}
	})

	t.Run("case=must-not-save-user-with-non-nullable-fields-as-nil", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		user := &models.User{}

		err := userRepository.Create(user)

		if err == nil {
			t.Errorf("Should not save with non nullable fields as nil")
		}
	})

	t.Run("case=must-not-save-user-with-repeated-username", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		user := &models.User{
			Username: "abc",
			Password: "111",
		}
		sameUsernameUser := &models.User{
			Username: "abc",
			Password: "888",
		}
		err := userRepository.Create(user)
		if err != nil {
			t.Errorf("Error saving user: %v", err)
		}

		err = userRepository.Create(sameUsernameUser)

		if err == nil {
			t.Errorf("Error saving user: Should not save with repeated username")
		}
	})

	t.Run("case=must-get-user-by-username", func(t *testing.T) {
		sqlite.DB.Exec("DELETE FROM users")
		username := "abcdef"
		user := &models.User{
			Username: username,
			Password: "abc",
		}
		userRepository.Create(user)

		savedUser, err := userRepository.GetByUsername(username)

		if err != nil {
			t.Errorf("Error fetching user: %v", err)
		}
		if savedUser.ID != user.ID {
			t.Errorf("Must fetch the same user")
		}
	})
}
