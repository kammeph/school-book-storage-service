package storage

import (
	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
)

func NewStorageRepository(store common.Store, serializer common.Serializer, broker common.MessageBroker) *common.Repository {
	domainEvents := []domain.Event{
		storage.StorageCreated{},
		storage.StorageRemoved{},
		storage.StorageNameSet{},
		storage.StorageLocationSet{},
	}
	for _, event := range domainEvents {
		serializer.Bind(event)
	}
	if broker != nil {
		broker.Subscribe(storage.StorageCreated{}, NewStorageCreatedEventHandler())
		broker.Subscribe(storage.StorageNameSet{}, NewStorageNameSetEventHandler())
	}
	return common.NewRepository(&storage.SchoolAggregateRoot{}, store, serializer, broker)
}
