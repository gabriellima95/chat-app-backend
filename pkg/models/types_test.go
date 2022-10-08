package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenericChat(t *testing.T) {
	var getNameTests = []struct {
		IsGroup      bool
		Username1    string
		UserID1      string
		Username2    string
		UserID2      string
		LoggedUserID string
		GroupName    string
		want         string
		desc         string
	}{
		{true, "thais", "030d8df7-0a89-46b4-9e34-1086af444c72", "gabriel", "ca5ebd01-2ca4-493f-980c-fac0b497901f", "ca5ebd01-2ca4-493f-980c-fac0b497901f", "amores", "amores", "GetName-should-return-group-name-when-is-group"},
		{false, "thais", "030d8df7-0a89-46b4-9e34-1086af444c72", "gabriel", "ca5ebd01-2ca4-493f-980c-fac0b497901f", "ca5ebd01-2ca4-493f-980c-fac0b497901f", "amores", "thais", "GetName-should-return-username-when-is-not-group-and-not-logged-user"},
		{false, "thais", "030d8df7-0a89-46b4-9e34-1086af444c72", "gabriel", "ca5ebd01-2ca4-493f-980c-fac0b497901f", "030d8df7-0a89-46b4-9e34-1086af444c72", "amores", "gabriel", "GetName-should-return-username-when-is-not-group-and-logged-user"},
	}

	for _, tt := range getNameTests {
		t.Run(tt.desc, func(t *testing.T) {
			user1 := User{
				ID:       uuid.MustParse(tt.UserID1),
				Username: tt.Username1,
				Password: "111",
			}
			user2 := User{
				ID:       uuid.MustParse(tt.UserID2),
				Username: tt.Username2,
				Password: "222",
			}
			chat := &GenericChat{
				Name:          tt.GroupName,
				LastMessageAt: time.Now(),
				LastMessage:   "oie",
				IsGroup:       tt.IsGroup,
				Users:         []User{user1, user2},
			}

			ans := chat.GetName(uuid.MustParse(tt.LoggedUserID))

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}

	userID := uuid.New()
	// user2ID := uuid.New()
	var getLastMessageTests = []struct {
		GenericChat  GenericChat
		User         User
		LoggedUserID uuid.UUID
		want         string
		desc         string
	}{
		{GenericChat{IsGroup: false, LastMessage: "oie"}, User{}, userID, "oie", "GetLastMessage-should-return-only-message-content-when-is-not-group"},
		{GenericChat{IsGroup: true, LastMessage: "oie", LastSenderID: userID}, User{}, userID, "Eu: oie", "GetLastMessage-should-return-'Eu: content'-when-is-group-and-is-the-logged-user"},
		{GenericChat{IsGroup: true, LastMessage: "oie", LastSenderID: userID}, User{ID: userID, Username: "Josias"}, uuid.New(), "Josias: oie", "GetLastMessage-should-return-'Username: content'-when-is-group-and-is-not-the-logged-user"},
	}

	for _, tt := range getLastMessageTests {
		t.Run(tt.desc, func(t *testing.T) {
			chat := tt.GenericChat
			chat.Users = []User{tt.User}
			ans := chat.GetLastMessage(tt.LoggedUserID)

			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
