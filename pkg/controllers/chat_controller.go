package controllers

import (
	"encoding/json"
	"msn/storage"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ChatController struct {
	chatRepository storage.ChatRepository
}

func NewChatController(chatRepository storage.ChatRepository) ChatController {
	return ChatController{
		chatRepository: chatRepository,
	}
}

func (c ChatController) ListChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	ID := params["user_id"]

	userID, err := uuid.Parse(ID)
	if err != nil {
		http.Error(w, "userID should be valid uuid", http.StatusBadRequest)
		return
	}

	chats, err := c.chatRepository.ListByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var chatListResponse []ChatResponseSchema
	for _, chat := range chats {
		c := ChatResponseSchema{
			ID:            chat.ID.String(),
			LastMessage:   chat.LastMessage,
			LastMessageAt: chat.LastMessageAt,
		}
		if userID == chat.User1ID {
			c.User = UserResponseSchema{
				ID:       chat.User1.ID.String(),
				Username: chat.User1.Username,
			}
			c.Contact = UserResponseSchema{
				ID:       chat.User2.ID.String(),
				Username: chat.User2.Username,
			}
		} else {
			c.User = UserResponseSchema{
				ID:       chat.User2.ID.String(),
				Username: chat.User2.Username,
			}
			c.Contact = UserResponseSchema{
				ID:       chat.User1.ID.String(),
				Username: chat.User1.Username,
			}

		}
		chatListResponse = append(chatListResponse, c)
	}

	json.NewEncoder(w).Encode(chatListResponse)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
