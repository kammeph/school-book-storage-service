package common

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Repository struct {
	prototype      reflect.Type
	store          Store
	serializer     Serializer
	commandHandler CommandHandler
	observers      []func(common.Event)
}

func NewRepository(prototype common.Aggregate, store Store, serializer Serializer, commandHandler CommandHandler) *Repository {
	t := reflect.TypeOf(prototype)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return &Repository{prototype: t, store: store, serializer: serializer, commandHandler: commandHandler}
}

func (r *Repository) NewAggregate() common.Aggregate {
	return reflect.New(r.prototype).Interface().(common.Aggregate)
}

func (r *Repository) Load(ctx context.Context, aggregateID uuid.UUID) (common.Aggregate, error) {
	records, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		return nil, err
	}

	aggregate := r.NewAggregate()
	if records == nil || len(records) == 0 {
		return aggregate, nil
	}
	for _, record := range records {
		event, err := r.serializer.UnmarshalEvent(record)
		if err != nil {
			return nil, err
		}

		err = aggregate.On(event)
		if err != nil {
			return nil, err
		}
	}
	return aggregate, nil
}

func (r *Repository) Save(ctx context.Context, command Command) error {
	aggregate, err := r.Load(ctx, command.AggregateID())
	if err != nil {
		return nil
	}
	if aggregate == nil {
		aggregate = r.NewAggregate()
	}
	domainEvents, err := r.commandHandler.Apply(ctx, aggregate, command)
	if err != nil {
		return nil
	}

	aggregateID := domainEvents[0].AggregateID()
	var records []Record
	for _, event := range domainEvents {
		record, err := r.serializer.MarshalEvent(event)
		if err != nil {
			return nil
		}
		records = append(records, record)

		// must work asynchronous
		for _, observer := range r.observers {
			observer(event)
		}
	}

	return r.store.Save(ctx, aggregateID, records...)
}
