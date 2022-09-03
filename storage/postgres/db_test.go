package postgres

import (
	"testing"
)

func TestChatRepository(t *testing.T) {
	Testing = true

	t.Run("case=must-connect-to-DB-and-run-migrations", func(t *testing.T) {
		SetupDatabase()
	})
}
