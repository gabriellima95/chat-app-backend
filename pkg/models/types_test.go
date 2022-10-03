package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenericChat(t *testing.T) {
	var tests = []struct {
		IsGroup   bool
		Username  string
		GroupName string
		want      string
		desc      string
	}{
		{true, "josias", "mazelados", "mazelados", "GetName-should-return-group-name-when-is-group"},
		{false, "josias", "mazelados", "josias", "GetName-should-return-username-when-is-not-group"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			user1 := User{
				ID:       uuid.New(),
				Username: "111",
				Password: "111",
			}
			user2 := User{
				ID:       uuid.New(),
				Username: tt.Username,
				Password: "222",
			}
			chat := &GenericChat{
				Name:          tt.GroupName,
				LastMessageAt: time.Now(),
				LastMessage:   "oie",
				IsGroup:       tt.IsGroup,
				Users:         []User{user1, user2},
			}

			ans := chat.GetName(user1.ID)

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
