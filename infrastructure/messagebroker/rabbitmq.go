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

func NewRabbitMQConnection() (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	return amqp.Dial(url)
}

type RabbitMQ struct {
	channel        *amqp.Channel
	exchange       string
	handlerByEvent map[string][]application.EventHandler
}

func NewRabbitMQ(connection *amqp.Connection, exchange string) (*RabbitMQ, error) {
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

func (broker *RabbitMQ) Subscribe(event domain.Event, handler application.EventHandler) error {
	eventName, _ := domain.EventType(event)
	handlers, _ := broker.handlerByEvent[eventName]
	handlers = append(handlers, handler)
	broker.handlerByEvent[eventName] = handlers
	return nil
}
