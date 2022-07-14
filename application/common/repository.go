package common

import (
	"context"
	"reflect"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Repository struct {
	prototype  reflect.Type
	store      Store
	serializer Serializer
	broker     MessageBroker
}

func NewRepository(prototype common.Aggregate, store Store, serializer Serializer, broker MessageBroker) *Repository {
	t := reflect.TypeOf(prototype)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return &Repository{prototype: t, store: store, serializer: serializer, broker: broker}
}

func (r *Repository) NewAggregate(id string) common.Aggregate {
	aggregate := reflect.New(r.prototype).Interface().(common.Aggregate)
	aggregate.SetAggregateID(id)
	return aggregate
}

func (r *Repository) Load(ctx context.Context, aggregateID string) (common.Aggregate, error) {
	records, err := r.store.Load(ctx, aggregateID)
	if err != nil {
		return nil, err
	}

	aggregate := r.NewAggregate(aggregateID)
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

func (r *Repository) Save(ctx context.Context, aggregate common.Aggregate) error {
	for _, event := range aggregate.DomainEvents() {
		record, err := r.serializer.MarshalEvent(event)
		if err != nil {
			return err
		}

		if err = r.store.Save(ctx, aggregate.AggregateID(), record); err != nil {
			return err
		}

		if r.broker != nil {
			r.broker.Publish(event)
		}
	}
	return nil
}
