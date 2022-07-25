package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain"
)

type MemoryRepository struct {
	storages []domain.StorageWithBooks
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{storages: []domain.StorageWithBooks{}}
}

func NewMemoryRepositoryWithStorages(storages []domain.StorageWithBooks) *MemoryRepository {
	return &MemoryRepository{storages: storages}
}

func (r *MemoryRepository) GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]domain.StorageWithBooks, error) {
	if r.storages == nil {
		return nil, errors.New("repository is not initialized")
	}
	storages := []domain.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	return storages, nil
}

func (r *MemoryRepository) GetStorageByID(ctx context.Context, schoolID, storageID string) (domain.StorageWithBooks, error) {
	if r.storages == nil {
		return domain.StorageWithBooks{}, errors.New("repository is not initialized")
	}
	storages := []domain.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.StorageID == storageID && storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	if len(storages) > 1 {
		return domain.StorageWithBooks{}, fmt.Errorf("more than one storage with ID %s found", storageID)
	}
	if len(storages) < 1 {
		return domain.StorageWithBooks{}, fmt.Errorf("no storage with ID %s found", storageID)
	}
	return storages[0], nil
}

func (r *MemoryRepository) GetStorageByName(ctx context.Context, schoolID, name string) (domain.StorageWithBooks, error) {
	if r.storages == nil {
		return domain.StorageWithBooks{}, errors.New("repository is not initialized")
	}
	storages := []domain.StorageWithBooks{}
	for _, storage := range r.storages {
		if storage.Name == name && storage.SchoolID == schoolID {
			storages = append(storages, storage)
		}
	}
	if len(storages) > 1 {
		return domain.StorageWithBooks{}, fmt.Errorf("more than one storage with name %s found", name)
	}
	if len(storages) < 1 {
		return domain.StorageWithBooks{}, fmt.Errorf("no storage with name %s found", name)
	}
	return storages[0], nil
}

func (r *MemoryRepository) InsertStorage(ctx context.Context, storage domain.StorageWithBooks) error {
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
