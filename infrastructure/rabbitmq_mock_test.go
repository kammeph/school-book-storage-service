package infrastructure_test

import (
	"errors"

	"github.com/kammeph/school-book-storage-service/infrastructure"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MockConntection struct {
	channelError bool
	closeError   bool
	channel      MockChannel
}

func NewMockConnection(
	channelError,
	exchangeDeclareError,
	queueDeclareError,
	queueBindError,
	consumeError,
	publishError,
	closeError bool,
) MockConntection {
	return MockConntection{
		channelError: channelError,
		closeError:   closeError,
		channel: MockChannel{
			exchangeDeclareError,
			queueDeclareError,
			queueBindError,
			consumeError,
			publishError,
		},
	}
}

type MockChannel struct {
	exchangeDeclareError bool
	queueDeclareError    bool
	queueBindError       bool
	consumeError         bool
	publishError         bool
}

var (
	errChannel         = errors.New("error while getting channel")
	errClose           = errors.New("error while closing connection")
	errExchangeDeclare = errors.New("error while declaring exchange")
	errQueueDeclare    = errors.New("error while declaring queue")
	errQueueBind       = errors.New("error while binding queue")
	errConsume         = errors.New("error while consuming from queue")
	errPublish         = errors.New("error while publishing message")
)

func (c MockConntection) Close() error {
	if c.closeError {
		return errClose
	}
	return nil
}

func (c MockConntection) Channel() (infrastructure.AmqpChannel, error) {
	if c.channelError {
		return nil, errChannel
	}
	return &c.channel, nil
}

func (ch *MockChannel) Close() error {
	return nil
}

func (ch *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	if ch.exchangeDeclareError {
		return errExchangeDeclare
	}
	return nil
}

func (ch *MockChannel) QueueDeclare(name string, durable, autoDelete, exclusiv, noWait bool, args amqp.Table) (amqp.Queue, error) {
	if ch.queueDeclareError {
		return amqp.Queue{}, errQueueDeclare
	}
	return amqp.Queue{Name: name}, nil
}

func (ch *MockChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	if ch.queueBindError {
		return errQueueBind
	}
	return nil
}

func (ch *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	if ch.consumeError {
		return nil, errConsume
	}
	return nil, nil
}

func (ch *MockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if ch.publishError {
		return errPublish
	}
	return nil
}
