package storage

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageWithBooksRepository interface {
	GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]storage.StorageWithBooks, error)
	GetStorageByID(ctx context.Context, schoolID, storageID string) (storage.StorageWithBooks, error)
	GetStorageByName(ctx context.Context, schoolID, name string) (storage.StorageWithBooks, error)
	InsertStorage(ctx context.Context, storage storage.StorageWithBooks) error
	DeleteStorage(ctx context.Context, storageID string) error
	UpdateStorageName(ctx context.Context, storageID, name string) error
	UpdateStorageLocation(ctx context.Context, storageID, location string) error
}
