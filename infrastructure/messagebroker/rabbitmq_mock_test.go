package messagebroker_test

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Conntection struct {
}

type Channel struct {
}

func (c *Conntection) Close() error {
	return nil
}

func (c *Conntection) Channel() (*Channel, error) {
	return &Channel{}, nil
}

func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return nil
}

func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusiv, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{Name: name}, nil
}

func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return nil
}

func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}
