package application

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
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
	QueryModel
}

func NewGetAllStorages(aggregateID string) GetAllStorages {
	return GetAllStorages{QueryModel: QueryModel{ID: aggregateID}}
}

type GetAllStoragesQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetAllStoragesQueryHandler(repository StorageWithBooksRepository) GetAllStoragesQueryHandler {
	return GetAllStoragesQueryHandler{repository: repository}
}

func (h GetAllStoragesQueryHandler) Handle(ctx context.Context, query GetAllStorages) ([]domain.StorageWithBooks, error) {
	return h.repository.GetAllStoragesBySchoolID(ctx, query.AggregateID())
}

type GetStorageByID struct {
	QueryModel
	StorageID string
}

func NewGetStorageByID(aggregateID, storageID string) GetStorageByID {
	return GetStorageByID{QueryModel: QueryModel{ID: aggregateID}, StorageID: storageID}
}

type GetStorageByIDQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetStorageByIDQueryHandler(repository StorageWithBooksRepository) GetStorageByIDQueryHandler {
	return GetStorageByIDQueryHandler{repository: repository}
}

func (h GetStorageByIDQueryHandler) Handle(ctx context.Context, query GetStorageByID) (domain.StorageWithBooks, error) {
	return h.repository.GetStorageByID(ctx, query.AggregateID(), query.StorageID)
}

type GetStorageByName struct {
	QueryModel
	Name string
}

func NewGetStorageByName(aggregateID string, name string) GetStorageByName {
	return GetStorageByName{QueryModel: QueryModel{ID: aggregateID}, Name: name}
}

type GetStorageByNameQueryHandler struct {
	repository StorageWithBooksRepository
}

func NewGetStorageByNameQueryHandler(repositoty StorageWithBooksRepository) GetStorageByNameQueryHandler {
	return GetStorageByNameQueryHandler{repository: repositoty}
}

func (h GetStorageByNameQueryHandler) Handle(ctx context.Context, query GetStorageByName) (domain.StorageWithBooks, error) {
	return h.repository.GetStorageByName(ctx, query.AggregateID(), query.Name)
}
