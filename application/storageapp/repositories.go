package storageapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
)

type StorageWithBooksRepository interface {
	GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]storagedomain.StorageWithBooks, error)
	GetStorageByID(ctx context.Context, schoolID, storageID string) (storagedomain.StorageWithBooks, error)
	GetStorageByName(ctx context.Context, schoolID, name string) (storagedomain.StorageWithBooks, error)
	InsertStorage(ctx context.Context, storage storagedomain.StorageWithBooks) error
	DeleteStorage(ctx context.Context, storageID string) error
	UpdateStorageName(ctx context.Context, storageID, name string) error
	UpdateStorageLocation(ctx context.Context, storageID, location string) error
}
