package application

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
)

type StorageWithBooksRepository interface {
	GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]domain.StorageWithBooks, error)
	GetStorageByID(ctx context.Context, schoolID, storageID string) (domain.StorageWithBooks, error)
	GetStorageByName(ctx context.Context, schoolID, name string) (domain.StorageWithBooks, error)
	InsertStorage(ctx context.Context, storage domain.StorageWithBooks) error
	DeleteStorage(ctx context.Context, storageID string) error
	UpdateStorageName(ctx context.Context, storageID, name string) error
	UpdateStorageLocation(ctx context.Context, storageID, location string) error
}
