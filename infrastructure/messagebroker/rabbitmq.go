package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"

	application "github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitUser     = "DB_USER"
	rabbitPassword = "DB_PASSWORD"
	rabbitHost     = "DB_HOST"
	rabbitPort     = "DB_PORT"
)

var (
	user     = utils.GetenvOrFallback(rabbitUser, "guest")
	password = utils.GetenvOrFallback(rabbitPassword, "guest")
	host     = utils.GetenvOrFallback(rabbitHost, "localhost")
	port     = utils.GetenvOrFallback(rabbitPort, "5672")
)

type AmqpConnection interface {
	Channel() (AmqpChannel, error)
	Close() error
}

type AmqpConnectionWrapper struct {
	connection *amqp.Connection
}

func (c AmqpConnectionWrapper) Close() error {
	return c.connection.Close()
}

func (c AmqpConnectionWrapper) Channel() (AmqpChannel, error) {
	return c.connection.Channel()
}

type AmqpChannel interface {
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusiv, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

func NewRabbitMQConnection() (AmqpConnection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return AmqpConnectionWrapper{}, err
	}
	return AmqpConnectionWrapper{conn}, nil
}

type RabbitMQ struct {
	channel        AmqpChannel
	exchange       string
	handlerByEvent map[string][]application.EventHandler
}

func NewRabbitMQ(connection AmqpConnection, exchange string) (*RabbitMQ, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	broker := &RabbitMQ{
		channel:        channel,
		exchange:       exchange,
		handlerByEvent: map[string][]application.EventHandler{},
	}

	err = channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, err
	}
	err = channel.QueueBind(q.Name, "", exchange, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range msgs {
			eventName := msg.RoutingKey
			handlers, ok := broker.handlerByEvent[eventName]
			if !ok {
				return
			}
			for _, handler := range handlers {
				handler.Handle(context.Background(), msg.Body)
			}
		}
	}()

	return broker, nil
}

func (broker *RabbitMQ) Publish(event domain.Event) error {
	eventType, _ := domain.EventType(event)
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{Body: eventData}
	return broker.channel.Publish(broker.exchange, eventType, false, false, msg)
}

func (broker *RabbitMQ) Subscribe(event domain.Event, handler application.EventHandler) {
	eventName, _ := domain.EventType(event)
	handlers, _ := broker.handlerByEvent[eventName]
	handlers = append(handlers, handler)
	broker.handlerByEvent[eventName] = handlers
}
