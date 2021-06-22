package rabbitmq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go-hackathon/src/common/application/messaging"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/common/infrastructure"
)

const queuePostfix = "_queue"
const exchangePostfix = "_exchange"
const routingPostfix = "_routing"

type rabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  cmd.AMQPConfig
}

func NewProducer(c cmd.AMQPConfig) messaging.AMQPProducer {
	client := newClient(c)

	_, err := client.channel.QueueDeclare(client.queueName(), false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = client.channel.ExchangeDeclare(client.exchangeName(), client.config.ExchangeType, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = client.channel.QueueBind(client.queueName(), client.routingName(), client.exchangeName(), false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NewConsumer(c cmd.AMQPConfig) messaging.AMQPConsumer {
	client := newClient(c)

	_, err := client.channel.QueueDeclare(client.queueName(), false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func newClient(c cmd.AMQPConfig) *rabbitMQClient {
	conn, err := amqp.Dial(c.ServerUrl)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &rabbitMQClient{conn, channel, c}
}

func (r *rabbitMQClient) Close() {
	infrastructure.Close(r.channel, "rabbit mq channel")
	infrastructure.Close(r.conn, "rabbit mq connection")
}

func (r *rabbitMQClient) Publish(msg string) error {
	return r.channel.Publish(r.exchangeName(), r.routingName(), false, false, amqp.Publishing{ContentType: "text/json", Body: []byte(msg)})
}

func (r *rabbitMQClient) Consume(consumer func(msg string) error) {
	messages, err := r.channel.Consume(r.queueName(), "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		err := consumer(string(message.Body))
		if err == nil {
			err = message.Ack(false)
			if err != nil {
				log.Error(err)
			}
		} else {
			err = message.Nack(false, true)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (r *rabbitMQClient) queueName() string {
	return r.config.QueueName + queuePostfix
}

func (r *rabbitMQClient) exchangeName() string {
	return r.config.QueueName + exchangePostfix
}

func (r *rabbitMQClient) routingName() string {
	return r.config.QueueName + routingPostfix
}
