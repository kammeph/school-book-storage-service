package storage

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageWithBooksRepository interface {
	GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]storage.StorageWithBooks, error)
	GetStorageByID(ctx context.Context, storageID string) (storage.StorageWithBooks, error)
	GetStorageByName(ctx context.Context, name string) (storage.StorageWithBooks, error)
	InsertStorage(ctx context.Context, storage storage.StorageWithBooks)
	DeleteStorage(ctx context.Context, storageID string)
	UpdateStorageName(ctx context.Context, storageID, name string)
	UpdateStorageLocation(ctx context.Context, storageID, location string)
}
