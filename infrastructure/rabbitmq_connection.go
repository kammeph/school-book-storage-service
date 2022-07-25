package infrastructure

import (
	"fmt"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbituser     = utils.GetenvOrFallback("RABBIT_USER", "guest")
	rabbitpassword = utils.GetenvOrFallback("RABBIT_PASSWORD", "guest")
	rabbithost     = utils.GetenvOrFallback("RABBIT_HOST", "localhost")
	rabbitport     = utils.GetenvOrFallback("RABBIT_PORT", "5672")
)

func NewRabbitMQConnection() AmqpConnection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbituser, rabbitpassword, rabbithost, rabbitport)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to rabbit mq.")
	return AmqpConnectionWrapper{conn}
}

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
