package pubsub

type PublisherMock struct {
	// UploadFn       func(filename string, content []byte) (string, error)
	// UploadBase64Fn func(filename string, content string) (string, error)
}

func (f *PublisherMock) Publish(msg Message) error {
	return nil
}
