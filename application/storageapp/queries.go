package storageapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
)

type StorageQueryHandlers struct {
	GetAllHandler           GetAllStoragesQueryHandler
	GetStorageByIDHandler   GetStorageByIDQueryHandler
	GetStorageByNameHandler GetStorageByNameQueryHandler
}

func NewStorageQueryHandlers(repository StorageWithBooksRepository) StorageQueryHandlers {
	return StorageQueryHandlers{
		GetAllHandler:           NewGetAllStoragesQueryHandler(repository),
		GetStorageByIDHandler:   NewGetStorageByIDQueryHandler(repository),
		GetStorageByNameHandler: NewGetStorageByNameQueryHandler(repository),
	}
}

type GetAllStorages struct {
	application.QueryModel
}

func NewGetAllStorages(aggregateID string) GetAllStorages {
	return GetAllStorages{QueryModel: application.QueryModel{ID: aggregateID}}
}

type GetAllStoragesQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetAllStoragesQueryHandler(repository StorageWithBooksRepository) GetAllStoragesQueryHandler {
	return GetAllStoragesQueryHandler{repository: repository}
}

func (h GetAllStoragesQueryHandler) Handle(ctx context.Context, query GetAllStorages) ([]storagedomain.StorageWithBooks, error) {
	return h.repository.GetAllStoragesBySchoolID(ctx, query.AggregateID())
}

type GetStorageByID struct {
	application.QueryModel
	StorageID string
}

func NewGetStorageByID(aggregateID, storageID string) GetStorageByID {
	return GetStorageByID{QueryModel: application.QueryModel{ID: aggregateID}, StorageID: storageID}
}

type GetStorageByIDQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetStorageByIDQueryHandler(repository StorageWithBooksRepository) GetStorageByIDQueryHandler {
	return GetStorageByIDQueryHandler{repository: repository}
}

func (h GetStorageByIDQueryHandler) Handle(ctx context.Context, query GetStorageByID) (storagedomain.StorageWithBooks, error) {
	return h.repository.GetStorageByID(ctx, query.AggregateID(), query.StorageID)
}

type GetStorageByName struct {
	application.QueryModel
	Name string
}

func NewGetStorageByName(aggregateID string, name string) GetStorageByName {
	return GetStorageByName{QueryModel: application.QueryModel{ID: aggregateID}, Name: name}
}

type GetStorageByNameQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetStorageByNameQueryHandler(repositoty StorageWithBooksRepository) GetStorageByNameQueryHandler {
	return GetStorageByNameQueryHandler{repository: repositoty}
}

func (h GetStorageByNameQueryHandler) Handle(ctx context.Context, query GetStorageByName) (storagedomain.StorageWithBooks, error) {
	return h.repository.GetStorageByName(ctx, query.AggregateID(), query.Name)
}
