package common

import (
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	Publish(event domain.Event) error
	Subscribe(event domain.Event, handler EventHandler) error
}

// Connection interface for the rabbit mq message broker
type AmqpConnection interface {
	Close() error
	Channel() (AmqpChannel, error)
}

// Channel interface for the rabbit mq message broker
type AmqpChannel interface {
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusiv, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}
