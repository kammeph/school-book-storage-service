package storage_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/storage"
	"github.com/kammeph/school-book-storage-service/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetAllStoragesBySchoolID(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		collection  string
		schoolID    string
		expectError bool
		err         error
	}{
		{
			name:        "get all storages by school id",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "school",
			expectError: false,
			err:         nil,
		},
		{
			name:        "get all storages by school id error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "error",
			expectError: true,
			err:         errors.New("mock-error"),
		},
		{
			name:        "get all storages by school id cursor error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "cursor-error",
			expectError: true,
			err:         errors.New("mock-cursor-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("Find", context.Background(), bson.D{{Key: "schoolId", Value: "school"}}).
				Return([]domain.StorageWithBooks{
					{
						SchoolID:  test.schoolID,
						StorageID: "storage1",
						Name:      "Closet 1",
						Location:  "Room 101",
					}},
					nil,
				)
			collection.(*mocks.MockCollection).
				On("Find", context.Background(), bson.D{{Key: "schoolId", Value: "error"}}).
				Return(nil, errors.New("mock-error"))
			collection.(*mocks.MockCollection).
				On("Find", context.Background(), bson.D{{Key: "schoolId", Value: "cursor-error"}}).
				Return(nil, nil)
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			storages, err := repository.GetAllStoragesBySchoolID(context.Background(), test.schoolID)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				assert.Nil(t, storages)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, storages, 1)
		})
	}
}

func TestGetStoragesByID(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		collection  string
		schoolID    string
		storageID   string
		expectError bool
		err         error
	}{
		{
			name:        "get storage by id",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "school",
			storageID:   "storage1",
			expectError: false,
			err:         nil,
		},
		{
			name:        "get storage by id error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "error",
			storageID:   "error",
			expectError: true,
			err:         errors.New("mock-find-one-error"),
		},
		{
			name:        "get storage by id decode error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "decode-error",
			storageID:   "decode-error",
			expectError: true,
			err:         errors.New("mock-decode-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "storageId", Value: "storage1"}}).
				Return(&mocks.MockSingleResult{Storage: &domain.StorageWithBooks{
					SchoolID:  test.schoolID,
					StorageID: test.storageID,
					Name:      "Closet 1",
					Location:  "Room 101",
				}})
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "storageId", Value: "error"}}).
				Return(&mocks.MockSingleResult{Error: errors.New("mock-find-one-error")})
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "storageId", Value: "decode-error"}}).
				Return(&mocks.MockSingleResult{Storage: nil})
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			storage, err := repository.GetStorageByID(context.Background(), test.schoolID, test.storageID)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
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
		database    string
		collection  string
		schoolID    string
		storageName string
		expectError bool
		err         error
	}{
		{
			name:        "get storage by name",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "school",
			storageName: "Closet 1",
			expectError: false,
			err:         nil,
		},
		{
			name:        "get storage by name error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "error",
			storageName: "error",
			expectError: true,
			err:         errors.New("mock-find-one-error"),
		},
		{
			name:        "get all storages by name decode error",
			database:    "testdb",
			collection:  "testcollection",
			schoolID:    "decode-error",
			storageName: "decode-error",
			expectError: true,
			err:         errors.New("mock-decode-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "name", Value: "Closet 1"}}).
				Return(&mocks.MockSingleResult{Storage: &domain.StorageWithBooks{
					SchoolID:  test.schoolID,
					StorageID: test.storageName,
					Name:      "Closet 1",
					Location:  "Room 101",
				}})
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "name", Value: "error"}}).
				Return(&mocks.MockSingleResult{Error: errors.New("mock-find-one-error")})
			collection.(*mocks.MockCollection).
				On("FindOne", context.Background(), bson.D{{Key: "name", Value: "decode-error"}}).
				Return(&mocks.MockSingleResult{Storage: nil})
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			storage, err := repository.GetStorageByName(context.Background(), test.schoolID, test.storageName)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				assert.Zero(t, storage)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, storage)
		})
	}
}

func TestInsertStorage(t *testing.T) {
	tests := []struct {
		name            string
		database        string
		collection      string
		schoolID        string
		storageID       string
		storageName     string
		storageLocation string
		expectError     bool
		err             error
	}{
		{
			name:            "insert storage",
			database:        "testdb",
			collection:      "testcollection",
			schoolID:        "school",
			storageID:       "storage1",
			storageName:     "Closet 1",
			storageLocation: "Room 101",
			expectError:     false,
			err:             nil,
		},
		{
			name:            "insert storage error",
			database:        "testdb",
			collection:      "testcollection",
			schoolID:        "error",
			storageID:       "storage1",
			storageName:     "Closet 1",
			storageLocation: "Room 101",
			expectError:     true,
			err:             errors.New("mock-insert-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("InsertOne", context.Background(), bson.D{
					{Key: "storageId", Value: "storage1"},
					{Key: "schoolId", Value: "school"},
					{Key: "name", Value: "Closet 1"},
					{Key: "location", Value: "Room 101"},
					{Key: "books", Value: []domain.BookInStorage{}},
				}).
				Return(nil, nil)
			collection.(*mocks.MockCollection).
				On("InsertOne", context.Background(), bson.D{
					{Key: "storageId", Value: "storage1"},
					{Key: "schoolId", Value: "error"},
					{Key: "name", Value: "Closet 1"},
					{Key: "location", Value: "Room 101"},
					{Key: "books", Value: []domain.BookInStorage{}},
				}).
				Return(nil, errors.New("mock-insert-error"))
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			storage := domain.StorageWithBooks{
				StorageID: test.storageID,
				SchoolID:  test.schoolID,
				Name:      test.storageName,
				Location:  test.storageLocation,
				Books:     []domain.BookInStorage{},
			}
			err := repository.InsertStorage(context.Background(), storage)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestDeleteStorage(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		collection  string
		storageID   string
		expectError bool
		err         error
	}{
		{
			name:        "delete storage",
			database:    "testdb",
			collection:  "testcollection",
			storageID:   "storage1",
			expectError: false,
			err:         nil,
		},
		{
			name:        "delete storage error",
			database:    "testdb",
			collection:  "testcollection",
			storageID:   "error",
			expectError: true,
			err:         errors.New("mock-delete-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("DeleteOne", context.Background(), bson.D{{Key: "storageId", Value: "storage1"}}).
				Return(nil, nil)
			collection.(*mocks.MockCollection).
				On("DeleteOne", context.Background(), bson.D{{Key: "storageId", Value: "error"}}).
				Return(nil, errors.New("mock-delete-error"))
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			err := repository.DeleteStorage(context.Background(), test.storageID)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUpdateStorageName(t *testing.T) {
	tests := []struct {
		name        string
		database    string
		collection  string
		storageID   string
		storageName string
		expectError bool
		err         error
	}{
		{
			name:        "update storage name",
			database:    "testdb",
			collection:  "testcollection",
			storageID:   "storage1",
			storageName: "renamed",
			expectError: false,
			err:         nil,
		},
		{
			name:        "update storage name error",
			database:    "testdb",
			collection:  "testcollection",
			storageID:   "error",
			storageName: "error",
			expectError: true,
			err:         errors.New("mock-update-error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := mocks.NewMockClient()
			database := client.Database(test.database)
			collection := database.Collection(test.collection)
			collection.(*mocks.MockCollection).
				On("UpdateOne", context.Background(),
					bson.D{{Key: "storageId", Value: "storage1"}},
					bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "renamed"}}}}).
				Return(nil, nil)
			collection.(*mocks.MockCollection).
				On("UpdateOne", context.Background(),
					bson.D{{Key: "storageId", Value: "error"}},
					bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: "error"}}}}).
				Return(nil, errors.New("mock-update-error"))
			repository := storage.NewStorageWithBookRepository(client, test.database, test.collection)
			err := repository.UpdateStorageName(context.Background(), test.storageID, test.storageName)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
