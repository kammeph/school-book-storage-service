package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
	domain_common "github.com/kammeph/school-book-storage-service/domain/common"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/stretchr/testify/assert"
)

type EntityAggregate struct {
	domain_common.AggregateModel
}

func (e *EntityAggregate) On(event domain_common.Event) error {
	return nil
}

type memoryStore struct {
	eventsById map[string][]common.Record
}

func (s *memoryStore) Save(ctx context.Context, aggregateID string, records ...common.Record) error {
	if _, ok := s.eventsById[aggregateID]; !ok {
		s.eventsById[aggregateID] = []common.Record{}
	}
	history := append(s.eventsById[aggregateID], records...)
	s.eventsById[aggregateID] = history
	return nil
}

func (s *memoryStore) Load(ctx context.Context, aggregateID string) ([]common.Record, error) {
	_, ok := s.eventsById[aggregateID]
	if ok {
		return s.eventsById[aggregateID], nil
	}
	return nil, nil
}

var repository = common.NewRepository(
	&domain.SchoolAggregateRoot{},
	&memoryStore{eventsById: map[string][]common.Record{}},
	serializers.NewJSONSerializerWithEvents(
		domain.StorageCreated{},
		domain.StorageRemoved{},
		domain.StorageNameSet{},
		domain.StorageLocationSet{},
	),
	nil,
)

var entityRepository = common.NewRepository(
	&EntityAggregate{},
	&memoryStore{eventsById: map[string][]common.Record{}},
	serializers.NewJSONSerializerWithEvents(),
	nil,
)

func TestHandleAddStorage(t *testing.T) {
	handler := storage.NewAddStorageCommandHandler(repository)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	storageID, err := handler.Handle(ctx, add)
	assert.Nil(t, err)
	assert.NotZero(t, storageID, 3)
}

func TestHandleRemoveStorage(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(repository)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	dto, err := addHandler.Handle(ctx, add)
	removeHandler := storage.NewRemoveStorageCommandHandler(repository)
	remove := storage.RemoveStorage{CommandModel: common.CommandModel{ID: commandId}, StorageID: dto.StorageID, Reason: "test"}
	err = removeHandler.Handle(ctx, remove)
	assert.Nil(t, err)
}

func TestHandleSetStorageName(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(repository)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	dto, err := addHandler.Handle(ctx, add)
	setNameHandler := storage.NewSetStorageNameCommandHandler(repository)
	setName := storage.SetStorageName{
		CommandModel: common.CommandModel{ID: commandId},
		StorageID:    dto.StorageID,
		Name:         "storage name set",
		Reason:       "test",
	}
	err = setNameHandler.Handle(ctx, setName)
	assert.Nil(t, err)
}

func TestHandleSetStorageLocation(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(repository)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	dto, err := addHandler.Handle(ctx, add)
	setLocationHandler := storage.NewSetStorageLocationCommandHandler(repository)
	setLocation := storage.SetStorageLocation{
		CommandModel: common.CommandModel{ID: commandId},
		StorageID:    dto.StorageID,
		Location:     "location set",
		Reason:       "test",
	}
	err = setLocationHandler.Handle(ctx, setLocation)
	assert.Nil(t, err)
}

func TestIncorrectAggregateError(t *testing.T) {
	handler := storage.NewAddStorageCommandHandler(entityRepository)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	_, err := handler.Handle(ctx, add)
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.IncorrectAggregateTypeError(&EntityAggregate{}))
}
