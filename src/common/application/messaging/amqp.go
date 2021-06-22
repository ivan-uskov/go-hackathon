package messaging

type AMQPProducer interface {
	Publish(msg string) error
	Close()
}
