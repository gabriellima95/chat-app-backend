package localfiles

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	RootPath string
}

func NewLocalFileStorage(rootPath string) LocalFileStorage {
	return LocalFileStorage{RootPath: rootPath}
}

func (l LocalFileStorage) Upload(filename string, content []byte) (string, error) {
	filePath := filepath.Join(l.RootPath, filename)
	err := ioutil.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}
	return filePath, nil
}

func (l LocalFileStorage) UploadBase64(filename string, content string, chatID string) (string, error) {
	filePath := filepath.Join(l.RootPath, chatID, filename)
	dirPath := filepath.Join(l.RootPath, chatID)
	decoded, err := DecodeBase64(content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %w", err)
	}
	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("error creating directory: %s", err.Error())
	}

	err = os.WriteFile(filePath, decoded, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}
	return filePath, nil
}

func DecodeBase64(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}
