package localfiles_test

import (
	"msn/storage/localfiles"
	"os"
	"testing"
)

func TestLocalFileStorage_Upload(t *testing.T) {
	storage := localfiles.NewLocalFileStorage("")
	filepath, err := storage.Upload("test.txt", []byte("test"))
	if err != nil {
		t.Errorf("Upload() error = %v", err)
	}
	if filepath != "test.txt" {
		t.Errorf("Upload() filepath = %v", filepath)
	}

	//Read file created and compare content
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	if string(content) != "test" {
		t.Errorf("UploadBase64() content = %v", string(content))
	}

	//Delete file
	err = os.Remove(filepath)
	if err != nil {
		t.Errorf("Error deleting file: %v", err)
	}
}

func TestLocalFileStorage_UploadBase64(t *testing.T) {
	storage := localfiles.NewLocalFileStorage("")
	// string "dGVzdA==" is base64 encoded "test"
	filepath, err := storage.UploadBase64("test.txt", "dGVzdA==", "123")
	if err != nil {
		t.Errorf("UploadBase64() error = %v", err)
	}
	if filepath != "123/test.txt" {
		t.Errorf("UploadBase64() filepath = %v", filepath)
	}

	//Read file created and compare content
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	if string(content) != "test" {
		t.Errorf("UploadBase64() content = %v", string(content))
	}

	// Delete directory
	err = os.RemoveAll("123")
	if err != nil {
		t.Errorf("Error deleting direactory: %v", err)
	}
}
