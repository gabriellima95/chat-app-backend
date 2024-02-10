package goroutine

import "msn/pubsub"

type Publisher struct {
	broker Broker
}

func NewPublisher(broker Broker) Publisher {
	return Publisher{
		broker: broker,
	}
}

func (p Publisher) Publish(msg pubsub.Message) error {
	p.broker.ReceiveMessage(msg)
	return nil
}
