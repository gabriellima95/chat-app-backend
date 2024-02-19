package goroutine

import (
	"log"
	"msn/pubsub"
)

type subscription struct {
	channel    chan pubsub.Message
	subscriber pubsub.Subscriber
}

func (s subscription) Start() {
	log.Printf("Subscriber %s waiting messages", s.subscriber.Name())
	for {
		msg := <-s.channel
		log.Printf("Receiving message on subscriber: %s", s.subscriber.Name())

		err := s.subscriber.Run(msg)
		if err != nil {
			log.Printf("Error on subscription %s, retry policy not set", s.subscriber.Name())
		}
	}
}

func (s subscription) Foward(msg pubsub.Message) {
	s.channel <- msg
}

func (s subscription) Name() string {
	return s.subscriber.Name()
}
