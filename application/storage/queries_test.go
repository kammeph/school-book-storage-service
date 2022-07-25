package storage_test

import (
	"context"
	"testing"

	application "github.com/kammeph/school-book-storage-service/application/storage"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	infrastructure "github.com/kammeph/school-book-storage-service/infrastructure/storage"
	"github.com/stretchr/testify/assert"
)

var (
	storage1School1        = domain.NewStorageWithBooks("school1", "storage1School1", "Closet 1", "Room 101")
	storage2School1        = domain.NewStorageWithBooks("school1", "storage2School1", "Closet 2", "Room 101")
	storage1School2        = domain.NewStorageWithBooks("school2", "storage1School2", "Closet 1", "Room 203")
	emptyRepository        = infrastructure.NewMemoryRepository()
	repositoryWithStorages = infrastructure.NewMemoryRepositoryWithStorages(
		[]domain.StorageWithBooks{storage1School1, storage2School1, storage1School2})
)

func TestGetAllStorages(t *testing.T) {
	tests := []struct {
		name             string
		repository       application.StorageWithBooksRepository
		queryID          string
		numberOfStorages int
		expectError      bool
	}{
		{
			name:             "get from school 1",
			repository:       repositoryWithStorages,
			queryID:          "school1",
			numberOfStorages: 2,
			expectError:      false,
		},
		{
			name:             "get from school 2",
			repository:       repositoryWithStorages,
			queryID:          "school2",
			numberOfStorages: 1,
			expectError:      false,
		},
		{
			name:             "quer empty repository",
			repository:       emptyRepository,
			queryID:          "school2",
			numberOfStorages: 0,
			expectError:      false,
		},
		{
			name:             "not initialized repository",
			repository:       &infrastructure.MemoryRepository{},
			queryID:          "school1",
			numberOfStorages: 0,
			expectError:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := application.NewGetAllStoragesQueryHandler(test.repository)
			query := application.NewGetAllStorages(test.queryID)
			ctx := context.Background()
			storages, err := handler.Handle(ctx, query)
			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, storages)
			assert.Len(t, storages, test.numberOfStorages)

		})
	}
}
func TestGetStoragesByID(t *testing.T) {
	tests := []struct {
		name        string
		repository  application.StorageWithBooksRepository
		queryID     string
		storageID   string
		expectError bool
	}{
		{
			name:        "get storage 1 from school 1",
			repository:  repositoryWithStorages,
			queryID:     "school1",
			storageID:   "storage1School1",
			expectError: false,
		},
		{
			name:        "get storage 2 from school 1",
			repository:  repositoryWithStorages,
			queryID:     "school1",
			storageID:   "storage2School1",
			expectError: false,
		},
		{
			name:        "get storage 1 from school 2",
			repository:  repositoryWithStorages,
			queryID:     "school2",
			storageID:   "storage1School2",
			expectError: false,
		},
		{
			name:        "query empty repository",
			repository:  emptyRepository,
			queryID:     "school1",
			storageID:   "storage2School1",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := application.NewGetStorageByIDQueryHandler(test.repository)
			query := application.NewGetStorageByID(test.queryID, test.storageID)
			ctx := context.Background()
			storage, err := handler.Handle(ctx, query)
			if test.expectError {
				assert.Error(t, err)
				assert.Zero(t, storage)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, storage)
		})
	}
}

func TestGetStoragesByName(t *testing.T) {
	tests := []struct {
		name        string
		repository  application.StorageWithBooksRepository
		queryID     string
		storaNamege string
		expectError bool
	}{
		{
			name:        "school 1 closet 1",
			repository:  repositoryWithStorages,
			queryID:     "school1",
			storaNamege: "Closet 1",
			expectError: false,
		},
		{
			name:        "school 1 closet 2",
			repository:  repositoryWithStorages,
			queryID:     "school1",
			storaNamege: "Closet 2",
			expectError: false,
		},
		{
			name:        "school 2 closet 1",
			repository:  repositoryWithStorages,
			queryID:     "school2",
			storaNamege: "Closet 1",
			expectError: false,
		},
		{
			name:        "try closet 2 from school 2",
			repository:  repositoryWithStorages,
			queryID:     "school2",
			storaNamege: "Closet 2",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := application.NewGetStorageByNameQueryHandler(test.repository)
			query := application.NewGetStorageByName(test.queryID, test.storaNamege)
			ctx := context.Background()
			storage, err := handler.Handle(ctx, query)
			if test.expectError {
				assert.Error(t, err)
				assert.Zero(t, storage)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, storage)
		})
	}
}
