package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageQueryHandlers struct {
	GetAllHandler           GetAllStoragesQueryHandler
	GetStorageByIDHandler   GetStorageByIDQueryHandler
	GetStorageByNameHandler GetStorageByNameQueryHandler
}

func NewStorageQueryHandlers(repository *common.Repository) StorageQueryHandlers {
	return StorageQueryHandlers{
		GetAllHandler:           NewGetAllStoragesQueryHandler(repository),
		GetStorageByIDHandler:   NewGetStorageByIDQueryHandler(repository),
		GetStorageByNameHandler: NewGetStorageByNameQueryHandler(repository),
	}
}

type GetAllStorages struct {
	common.QueryModel
}

func NewGetAllStorages(aggregateID uuid.UUID) GetAllStorages {
	return GetAllStorages{QueryModel: common.QueryModel{ID: aggregateID}}
}

type GetAllStoragesQueryHandler struct {
	repository *common.Repository
}

func NewGetAllStoragesQueryHandler(repository *common.Repository) GetAllStoragesQueryHandler {
	return GetAllStoragesQueryHandler{repository: repository}
}

func (h GetAllStoragesQueryHandler) Handle(ctx context.Context, query GetAllStorages) ([]domain.Storage, error) {
	aggregate, err := h.repository.Load(ctx, query.AggregateID())
	if err != nil {
		return nil, err
	}
	schoolAggregate, ok := aggregate.(*domain.SchoolAggregateRoot)
	if !ok {
		return nil, IncorrectAggregateTypeError(aggregate)
	}
	if schoolAggregate.School.Storages == nil {
		schoolAggregate.School.Storages = []domain.Storage{}
	}
	return schoolAggregate.School.Storages, nil
}

type GetStorageByID struct {
	common.QueryModel
	StorageID uuid.UUID
}

func NewGetStorageByID(aggregateID, storageID uuid.UUID) GetStorageByID {
	return GetStorageByID{QueryModel: common.QueryModel{ID: aggregateID}, StorageID: storageID}
}

type GetStorageByIDQueryHandler struct {
	repository *common.Repository
}

func NewGetStorageByIDQueryHandler(repository *common.Repository) GetStorageByIDQueryHandler {
	return GetStorageByIDQueryHandler{repository: repository}
}

func (h GetStorageByIDQueryHandler) Handle(ctx context.Context, query GetStorageByID) (domain.Storage, error) {
	aggregate, err := h.repository.Load(ctx, query.AggregateID())
	if err != nil {
		return domain.Storage{}, err
	}
	schoolAggregate, ok := aggregate.(*domain.SchoolAggregateRoot)
	if !ok {
		return domain.Storage{}, IncorrectAggregateTypeError(aggregate)
	}
	storage, _, err := schoolAggregate.GetStorageByID(query.StorageID)
	if err != nil {
		return domain.Storage{}, err
	}
	return *storage, nil
}

type GetStorageByName struct {
	common.QueryModel
	Name string
}

func NewGetStorageByName(aggregateID uuid.UUID, name string) GetStorageByName {
	return GetStorageByName{QueryModel: common.QueryModel{ID: aggregateID}, Name: name}
}

type GetStorageByNameQueryHandler struct {
	repository *common.Repository
}

func NewGetStorageByNameQueryHandler(repositoty *common.Repository) GetStorageByNameQueryHandler {
	return GetStorageByNameQueryHandler{repository: repositoty}
}

func (h GetStorageByNameQueryHandler) Handle(ctx context.Context, query GetStorageByName) (domain.Storage, error) {
	aggregate, err := h.repository.Load(ctx, query.AggregateID())
	if err != nil {
		return domain.Storage{}, err
	}
	schoolAggregate, ok := aggregate.(*domain.SchoolAggregateRoot)
	if !ok {
		return domain.Storage{}, IncorrectAggregateTypeError(aggregate)
	}
	storage, err := schoolAggregate.GetStorageByName(query.Name)
	if err != nil {
		return domain.Storage{}, err
	}
	return *storage, nil
}
