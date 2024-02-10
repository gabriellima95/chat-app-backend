package goroutine

import (
	"log"
	"msn/pubsub"
)

type Broker struct {
	channel     chan pubsub.Message
	subscribers map[string][]pubsub.Subscriber
}

func NewBroker() Broker {
	topicChannel := make(chan pubsub.Message, 10)
	subscribers := make(map[string][]pubsub.Subscriber, 0)

	return Broker{
		channel:     topicChannel,
		subscribers: subscribers,
	}
}

func (b Broker) Subscribe(topic string, subscriber pubsub.Subscriber) {
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
}

func (b Broker) ReceiveMessage(msg pubsub.Message) {
	b.channel <- msg
	log.Println("Message Received")
}

func (b Broker) Broadcast() {
	log.Println("Broadcasting...")
	for {
		msg := <-b.channel
		log.Printf("Processing Message for topic %s", msg.Topic)
		subs, ok := b.subscribers[msg.Topic]
		if !ok {
			continue
		}
		for _, sub := range subs {
			log.Printf("Processing Message on subscriber %s", sub.Name())

			// tratar erro para fazer a retentativa
			sub.Run(msg)
		}
	}
}
