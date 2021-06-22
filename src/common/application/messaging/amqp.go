package messaging

type AMQPProducer interface {
	Publish(msg string) error
	Close()
}

type AMQPConsumer interface {
	Consume(consumer func(msg string) error)
	Close()
}
