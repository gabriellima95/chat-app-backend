package subscribers

import (
	"encoding/json"
	"log"
	"msn/pkg/models"
	"msn/pubsub"
	"msn/storage"

	"github.com/google/uuid"
)

type UploadAttachmentsSubscriber struct {
	fileStorage       storage.FileStorageClient
	messageRepository storage.MessageRepository
}

func NewUploadAttachmentsSubscriber(fileStorage storage.FileStorageClient, messageRepository storage.MessageRepository) UploadAttachmentsSubscriber {
	return UploadAttachmentsSubscriber{
		fileStorage:       fileStorage,
		messageRepository: messageRepository,
	}
}

type UploadAttachmentsRequest struct {
	ID          string       `json:"id"`
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
		filepath, err := s.fileStorage.UploadBase64(attachment.Filename, attachment.Content, request.ChatID)
		if err != nil {
			log.Printf("Error on fileStorage.UploadBase64: %s", err.Error())
		}

		messageID, err := uuid.Parse(request.ID)
		if err != nil {
			log.Printf("Error on uuid.Parse: %s", err.Error())
		}

		attachment := models.Attachment{
			Path:      filepath,
			MessageID: messageID,
		}
		err = s.messageRepository.SaveAttachment(&attachment)
		if err != nil {
			log.Printf("Error on messageRepository.SaveAttachment: %s", err.Error())
		}
	}
	return nil
}
