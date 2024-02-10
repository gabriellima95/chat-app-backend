package pubsub

type Subscriber interface {
	Run(msg Message) error
	Name() string
}

type Publisher interface {
	Publish(msg Message) error
}

type Message struct {
	Topic   string
	Payload map[string]interface{}
}
