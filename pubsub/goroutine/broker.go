package goroutine

import (
	"log"
	"msn/pubsub"
)

type Broker struct {
	channel       chan pubsub.Message
	subscriptions map[string][]subscription
}

func NewBroker() Broker {
	topicChannel := make(chan pubsub.Message, 10)
	subscriptions := make(map[string][]subscription, 0)

	return Broker{
		channel:       topicChannel,
		subscriptions: subscriptions,
	}
}

func (b Broker) Subscribe(topic string, subscriber pubsub.Subscriber) {
	subscriptionChannel := make(chan pubsub.Message, 10)
	sub := subscription{
		channel:    subscriptionChannel,
		subscriber: subscriber,
	}
	b.subscriptions[topic] = append(b.subscriptions[topic], sub)
	go sub.Start()
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
		subs, ok := b.subscriptions[msg.Topic]
		if !ok {
			continue
		}
		for _, sub := range subs {
			log.Printf("Fowarding Message to subscriber %s", sub.Name())
			sub.Foward(msg)
		}
	}
}
