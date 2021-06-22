package events

import (
	"database/sql"
	"go-hackathon/src/common/application/events"
	"go-hackathon/src/common/application/messaging"
	eventsImpl "go-hackathon/src/common/infrastructure/events"
)

type Producer interface {
	Produce() error
}

func NewProducer(db *sql.DB, amqpProducer messaging.AMQPProducer) Producer {
	return &producer{eventsImpl.NewStoredEventUnitOfWork(db), amqpProducer}
}

type producer struct {
	uow          events.StoredEventUnitOfWork
	amqpProducer messaging.AMQPProducer
}

func (p *producer) Produce() error {
	return events.NewProduceEventsCommandHandler(p.uow, p.amqpProducer).Handle()
}
