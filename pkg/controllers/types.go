package controllers

import "time"

type UserRequestSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponseSchema struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	// Token    string `json:"token"`
}

type LoginResponseSchema struct {
	Token string `json:"token"`
	// Password string `json:"password"`
}

// type ChatRequestSchema struct {
// 	Chatname string `json:"username"`
// 	Password string `json:"password"`
// }

type ChatResponseSchema struct {
	ID            string             `json:"id"`
	LastMessage   string             `json:"last_message"`
	LastMessageAt time.Time          `json:"last_message_at"`
	User          UserResponseSchema `json:"user"`
	Contact       UserResponseSchema `json:"contact"`
}

type GenericChatResponseSchema struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
}

type MessageRequestSchema struct {
	Content     string       `json:"content"`
	ChatID      string       `json:"chat_id"`
	SenderID    string       `json:"sender_id"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	Filename    string `json:"filename"`
}

type MessageResponseSchema struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
}
