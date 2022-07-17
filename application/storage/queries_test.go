package storage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
)

func prepareAggregate() (string, string, error) {
	handler := storage.NewAddStorageCommandHandler(store, nil)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	storageID, err := handler.Handle(ctx, add)
	return commandId, storageID, err
}

// func TestGetAllStorages(t *testing.T) {
// 	queryID, _, err := prepareAggregate()
// 	assert.Nil(t, err)
// 	getAllStorages := storage.GetAllStorages{QueryModel: common.QueryModel{ID: queryID}}
// 	handler := storage.NewGetAllStoragesQueryHandler(repository)
// 	ctx := context.Background()
// 	storages, err := handler.Handle(ctx, getAllStorages)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, storages)
// 	assert.Len(t, storages, 1)
// }

// func TestGetStoragesByID(t *testing.T) {
// 	queryID, storageID, err := prepareAggregate()
// 	assert.Nil(t, err)
// 	getStorageByID := storage.GetStorageByID{QueryModel: common.QueryModel{ID: queryID}, StorageID: storageID}
// 	handler := storage.NewGetStorageByIDQueryHandler(repository)
// 	ctx := context.Background()
// 	storage, err := handler.Handle(ctx, getStorageByID)
// 	assert.Nil(t, err)
// 	assert.NotZero(t, storage)
// }

// func TestGetStorageIDNotFound(t *testing.T) {
// 	queryID, _, err := prepareAggregate()
// 	storageID := uuid.New().String()
// 	assert.Nil(t, err)
// 	getStorageByID := storage.GetStorageByID{QueryModel: common.QueryModel{ID: queryID}, StorageID: storageID}
// 	handler := storage.NewGetStorageByIDQueryHandler(repository)
// 	ctx := context.Background()
// 	storage, err := handler.Handle(ctx, getStorageByID)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err, domain.StorageIDNotFoundError(storageID))
// 	assert.Zero(t, storage)
// }

// func TestGetStoragesByName(t *testing.T) {
// 	queryID, _, err := prepareAggregate()
// 	assert.Nil(t, err)
// 	getStorageByName := storage.GetStorageByName{QueryModel: common.QueryModel{ID: queryID}, Name: "storage"}
// 	handler := storage.NewGetStorageByNameQueryHandler(repository)
// 	ctx := context.Background()
// 	storage, err := handler.Handle(ctx, getStorageByName)
// 	assert.Nil(t, err)
// 	assert.NotZero(t, storage)
// }

// func TestGetStoragesByNameNotFoundError(t *testing.T) {
// 	queryID, _, err := prepareAggregate()
// 	assert.Nil(t, err)
// 	getStorageByName := storage.GetStorageByName{QueryModel: common.QueryModel{ID: queryID}, Name: "unknown"}
// 	handler := storage.NewGetStorageByNameQueryHandler(repository)
// 	ctx := context.Background()
// 	storage, err := handler.Handle(ctx, getStorageByName)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err, domain.StorageByNameNotFoundError("unknown"))
// 	assert.Zero(t, storage)
// }
