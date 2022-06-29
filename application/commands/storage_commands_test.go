package commands_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/commands"
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/domain/aggregates"
	"github.com/kammeph/school-book-storage-service/domain/events"
	"github.com/stretchr/testify/assert"
)

type UnknowCommand struct {
	common.CommandModel
}

func TestApplyAddStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	handler := commands.StorageCommandHandler{}
	commandId := uuid.New()
	ctx := context.Background()
	add := commands.AddStorage{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	createdEvents, err := handler.Apply(ctx, &aggregate, add)
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 3)
	created, createdOk := createdEvents[0].(events.StorageCreated)
	assert.True(t, createdOk)
	assert.Equal(t, created.AggregateID(), commandId)
	assert.Equal(t, created.EventVersion(), 1)
	assert.NotZero(t, created.EventAt())
	nameSet, nameSetOk := createdEvents[1].(events.StorageNameSet)
	assert.True(t, nameSetOk)
	assert.Equal(t, nameSet.AggregateID(), commandId)
	assert.Equal(t, nameSet.EventVersion(), 2)
	assert.NotZero(t, nameSet.EventAt())
	assert.Equal(t, nameSet.Name, "storage")
	assert.Equal(t, nameSet.Reason, "initial create")
	locationSet, locationSetOk := createdEvents[2].(events.StorageLocationSet)
	assert.True(t, locationSetOk)
	assert.Equal(t, locationSet.AggregateID(), commandId)
	assert.Equal(t, locationSet.EventVersion(), 3)
	assert.NotZero(t, locationSet.EventAt())
	assert.Equal(t, locationSet.Location, "location")
	assert.Equal(t, locationSet.Reason, "initial create")
}

func TestApplyRemoveStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	handler := commands.StorageCommandHandler{}
	ctx := context.Background()
	remove := commands.RemoveStorage{CommandModel: common.CommandModel{}, Reason: "test"}
	createdEvents, err := handler.Apply(ctx, &aggregate, remove)
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	removed, ok := createdEvents[0].(events.StorageRemoved)
	assert.True(t, ok)
	assert.Equal(t, removed.EventVersion(), 1)
	assert.NotZero(t, removed.EventAt())
	assert.Equal(t, removed.Reason, "test")
}

func TestApplySetStorageName(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	handler := commands.StorageCommandHandler{}
	ctx := context.Background()
	setName := commands.SetStorageName{CommandModel: common.CommandModel{}, Name: "storage", Reason: "test"}
	createdEvents, err := handler.Apply(ctx, &aggregate, setName)
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	nameSet, ok := createdEvents[0].(events.StorageNameSet)
	assert.True(t, ok)
	assert.Equal(t, nameSet.EventVersion(), 1)
	assert.NotZero(t, nameSet.EventAt())
	assert.Equal(t, nameSet.Name, "storage")
	assert.Equal(t, nameSet.Reason, "test")
}

func TestApplySetStorageLocation(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	handler := commands.StorageCommandHandler{}
	ctx := context.Background()
	setLocation := commands.SetStorageLocation{CommandModel: common.CommandModel{}, Location: "location", Reason: "test"}
	createdEvents, err := handler.Apply(ctx, &aggregate, setLocation)
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	locationSet, ok := createdEvents[0].(events.StorageLocationSet)
	assert.True(t, ok)
	assert.Equal(t, locationSet.EventVersion(), 1)
	assert.NotZero(t, locationSet.EventAt())
	assert.Equal(t, locationSet.Location, "location")
	assert.Equal(t, locationSet.Reason, "test")
}

func TestApplyUnknowCommand(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	handler := commands.StorageCommandHandler{}
	ctx := context.Background()
	setLocation := UnknowCommand{CommandModel: common.CommandModel{}}
	_, err := handler.Apply(ctx, &aggregate, setLocation)
	assert.NotNil(t, err)
}
