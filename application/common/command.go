package common

import (
	"context"
	"fmt"

	domain "github.com/kammeph/school-book-storage-service/domain/common"
)

type Command interface {
	AggregateID() string
}

type CommandModel struct {
	ID string `json:"aggregateId"`
}

func (c CommandModel) AggregateID() string {
	return c.ID
}

type EventPublisher interface {
	Publish(ctx context.Context, events []domain.Event) error
}

type CommandHandlerModel struct {
	store     Store
	publisher EventPublisher
}

func NewCommandHandlerModel(store Store, publisher EventPublisher) *CommandHandlerModel {
	return &CommandHandlerModel{store, publisher}
}

func (h *CommandHandlerModel) SaveAndPublish(ctx context.Context, aggregate domain.Aggregate) error {
	if err := h.store.Save(ctx, aggregate); err != nil {
		return err
	}
	if h.publisher == nil {
		return nil
	}
	if err := h.publisher.Publish(ctx, aggregate.DomainEvents()); err != nil {
		fmt.Printf("Error while publishing events: %s", err)
	}
	return nil
}

func (h *CommandHandlerModel) Store() Store {
	return h.store
}

func (h *CommandHandlerModel) Publisher() EventPublisher {
	return h.publisher
}
