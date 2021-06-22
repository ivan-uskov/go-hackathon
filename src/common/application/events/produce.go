package events

import (
	"go-hackathon/src/common/application/messaging"
	"go-hackathon/src/common/model/events"
)

type StoredEventUnitOfWork interface {
	Execute(func(events.StoredEventRepository) error) error
	Lock(string, func() error) error
}

type ProduceEventsCommandHandler interface {
	Handle() error
}

type produceEventsCommandHandler struct {
	uow          StoredEventUnitOfWork
	amqpProducer messaging.AMQPProducer
}

func NewProduceEventsCommandHandler(uow StoredEventUnitOfWork, amqpProducer messaging.AMQPProducer) ProduceEventsCommandHandler {
	return &produceEventsCommandHandler{uow: uow, amqpProducer: amqpProducer}
}

func (p *produceEventsCommandHandler) Handle() error {
	var ee []events.StoredEvent
	err := p.uow.Execute(func(repository events.StoredEventRepository) error {
		var err error
		ee, err = repository.GetNotPublishedEvents()
		return err
	})
	if err != nil {
		return err
	}

	for _, e := range ee {
		err = p.amqpProducer.Publish(e.Body)
		if err != nil {
			return err
		}

		err = p.uow.Execute(func(repository events.StoredEventRepository) error {
			e.Publish()
			return repository.Store(e)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
