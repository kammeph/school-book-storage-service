package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitEventPublisher struct {
	channel  AmqpChannel
	exchange string
}

func NewRabbitEventPublisher(connection AmqpConnection, exchange string) (common.EventPublisher, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	if err = channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &RabbitEventPublisher{channel, exchange}, nil
}

func (p *RabbitEventPublisher) Publish(ctx context.Context, events []domain.Event) error {
	for _, event := range events {
		eventBytes, err := json.Marshal(event)
		if err != nil {
			return err
		}
		msg := amqp.Publishing{Body: eventBytes}
		return p.channel.Publish(p.exchange, event.EventType(), false, false, msg)

	}
	return nil
}
