package storage

type FileStorageClient interface {
	Upload(filename string, content []byte) (string, error)
	UploadBase64(filename, content, chatID string) (string, error)
	GeneratePresignedURL(objectKey string, lifetimeSecs int64) (string, error)
}
