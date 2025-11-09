package events

import (
	"context"
	"ride-sharing/shared/messaging"
)

type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMq
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMq) *TripEventPublisher {
	return &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}
func (p *TripEventPublisher) PublishTripCreate(ctx context.Context) error {

	return p.rabbitmq.PublishMessage(ctx, "Hello", "Hello World")
}
