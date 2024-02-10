package subscribers

import (
	"encoding/json"
	"log"
	"msn/pubsub"
	"msn/storage"
)

type UploadAttachmentsSubscriber struct {
	fileStorage storage.FileStorageClient
}

func NewUploadAttachmentsSubscriber(fileStorage storage.FileStorageClient) UploadAttachmentsSubscriber {
	return UploadAttachmentsSubscriber{
		fileStorage: fileStorage,
	}
}

type UploadAttachmentsRequest struct {
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

func (s UploadAttachmentsSubscriber) Name() string {
	return "upload-attachments-subscriber"
}

func (s UploadAttachmentsSubscriber) Run(msg pubsub.Message) error {
	payloadBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		log.Printf("Unable to Marshal message payload: %s", err.Error())
	}

	var request UploadAttachmentsRequest
	err = json.Unmarshal(payloadBytes, &request)
	if err != nil {
		log.Printf("Unable to Unmarshal message payload: %s", err.Error())
	}

	for _, attachment := range request.Attachments {
		_, err = s.fileStorage.UploadBase64(attachment.Filename, attachment.Content, request.ChatID)
		if err != nil {
			log.Printf("Error on fileStorage.UploadBase64: %s", err.Error())
		}
	}
	return nil
}
