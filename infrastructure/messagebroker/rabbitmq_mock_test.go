package messagebroker_test

import (
	"errors"

	"github.com/kammeph/school-book-storage-service/infrastructure/messagebroker"
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
	channelError         = errors.New("Error while getting channel")
	closeError           = errors.New("Error while closing connection")
	exchangeDeclareError = errors.New("Error while declaring exchange")
	queueDeclareError    = errors.New("Error while declaring queue")
	queueBindError       = errors.New("Error while binding queue")
	consumeError         = errors.New("Error while consuming from queue")
	publishError         = errors.New("Error while publishing message")
)

func (c MockConntection) Close() error {
	if c.closeError {
		return closeError
	}
	return nil
}

func (c MockConntection) Channel() (messagebroker.AmqpChannel, error) {
	if c.channelError {
		return nil, channelError
	}
	return &c.channel, nil
}

func (ch *MockChannel) Close() error {
	return nil
}

func (ch *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	if ch.exchangeDeclareError {
		return exchangeDeclareError
	}
	return nil
}

func (ch *MockChannel) QueueDeclare(name string, durable, autoDelete, exclusiv, noWait bool, args amqp.Table) (amqp.Queue, error) {
	if ch.queueDeclareError {
		return amqp.Queue{}, queueDeclareError
	}
	return amqp.Queue{Name: name}, nil
}

func (ch *MockChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	if ch.queueBindError {
		return queueBindError
	}
	return nil
}

func (ch *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	if ch.consumeError {
		return nil, consumeError
	}
	return nil, nil
}

func (ch *MockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if ch.publishError {
		return publishError
	}
	return nil
}
