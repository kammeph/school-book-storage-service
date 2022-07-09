package storage

import (
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
)

func NewStorageRepository(store common.Store, serializer common.Serializer) *common.Repository {
	serializer.Bind(
		storage.StorageCreated{},
		storage.StorageRemoved{},
		storage.StorageNameSet{},
		storage.StorageLocationSet{},
	)
	return common.NewRepository(&storage.SchoolAggregateRoot{}, store, serializer)
}
