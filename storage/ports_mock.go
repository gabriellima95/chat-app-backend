package storage

type FileStorageClientMock struct {
	// UploadFn       func(filename string, content []byte) (string, error)
	// UploadBase64Fn func(filename string, content string) (string, error)
}

func (f *FileStorageClientMock) Upload(filename string, content []byte) (string, error) {
	return "", nil
}

func (f *FileStorageClientMock) UploadBase64(filename string, content string) (string, error) {
	return "", nil
}
