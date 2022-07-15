package storage

import (
	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
)

func NewStorageRepository(store common.Store, serializer common.Serializer, broker common.MessageBroker) *common.Repository {
	domainEvents := []domain.Event{
		storage.StorageAdded{},
		storage.StorageRemoved{},
		storage.StorageRenamed{},
		storage.StorageRelocated{},
	}
	for _, event := range domainEvents {
		serializer.Bind(event)
	}
	if broker != nil {
		broker.Subscribe(storage.StorageAdded{}, NewStorageAddedEventHandler())
		broker.Subscribe(storage.StorageRenamed{}, NewStorageRenamedSetEventHandler())
	}
	return common.NewRepository(&storage.StorageAggregateRoot{}, store, serializer, broker)
}
