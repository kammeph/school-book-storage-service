package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain/storage"
)

type MemoryRepository struct {
	storages []storage.StorageWithBooks
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{storages: []storage.StorageWithBooks{}}
}

func NewMemoryRepositoryWithStorages(storages []storage.StorageWithBooks) *MemoryRepository {
	return &MemoryRepository{storages: storages}
}

func (r *MemoryRepository) GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]storage.StorageWithBooks, error) {
	if r.storages == nil {
		return nil, errors.New("repository is not initialized")
	}
	storages := []storage.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	return storages, nil
}

func (r *MemoryRepository) GetStorageByID(ctx context.Context, schoolID, storageID string) (storage.StorageWithBooks, error) {
	if r.storages == nil {
		return storage.StorageWithBooks{}, errors.New("repository is not initialized")
	}
	storages := []storage.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.StorageID == storageID && storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	if len(storages) > 1 {
		return storage.StorageWithBooks{}, fmt.Errorf("more than one storage with ID %s found", storageID)
	}
	if len(storages) < 1 {
		return storage.StorageWithBooks{}, fmt.Errorf("no storage with ID %s found", storageID)
	}
	return storages[0], nil
}

func (r *MemoryRepository) GetStorageByName(ctx context.Context, schoolID, name string) (storage.StorageWithBooks, error) {
	if r.storages == nil {
		return storage.StorageWithBooks{}, errors.New("repository is not initialized")
	}
	storages := []storage.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.Name == name && storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	if len(storages) > 1 {
		return storage.StorageWithBooks{}, fmt.Errorf("more than one storage with name %s found", name)
	}
	if len(storages) < 1 {
		return storage.StorageWithBooks{}, fmt.Errorf("no storage with name %s found", name)
	}
	return storages[0], nil
}

func (r *MemoryRepository) InsertStorage(ctx context.Context, storage storage.StorageWithBooks) error {
	r.storages = append(r.storages, storage)
	return nil
}

func (r *MemoryRepository) DeleteStorage(ctx context.Context, storageID string) error {
	for idx, storage := range r.storages {
		if storage.StorageID == storageID {
			r.storages = append(r.storages[:idx], r.storages[idx+1:]...)
			return nil
		}
	}
	return nil
}

func (r *MemoryRepository) UpdateStorageName(ctx context.Context, storageID, name string) error {
	for idx, storage := range r.storages {
		if storage.StorageID == storageID {
			r.storages[idx].Name = name
			return nil
		}
	}
	return nil
}

func (r *MemoryRepository) UpdateStorageLocation(ctx context.Context, storageID, location string) error {
	for idx, storage := range r.storages {
		if storage.StorageID == storageID {
			r.storages[idx].Location = location
			return nil
		}
	}
	return nil
}
