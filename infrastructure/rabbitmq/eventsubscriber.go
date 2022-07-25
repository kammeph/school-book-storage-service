package rabbitmq

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application/common"
)

type Subscription struct {
	channel AmqpChannel
	handler common.EventHandler
}

func NewSubscription(channel AmqpChannel, handler common.EventHandler) *Subscription {
	return &Subscription{channel, handler}
}

func (s *Subscription) Consume(exchange string) error {
	q, err := s.channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}
	if err := s.channel.QueueBind(q.Name, "", exchange, false, nil); err != nil {
		return err
	}

	msgs, err := s.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			if s.handler != nil {
				s.handler.Handle(context.Background(), msg.Body)
			}
		}
	}()
	return nil
}

type RabbitEventSubscriber struct {
	channel       AmqpChannel
	subscriptions []*Subscription
}

func NewRabbitEventSubscriber(connection AmqpConnection) (common.EventSubscriber, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitEventSubscriber{channel, []*Subscription{}}, nil
}

func (s *RabbitEventSubscriber) Subscribe(exchange string, handler common.EventHandler) error {
	subscription := NewSubscription(s.channel, handler)
	if err := subscription.Consume(exchange); err != nil {
		return err
	}
	s.subscriptions = append(s.subscriptions, subscription)
	return nil
}
