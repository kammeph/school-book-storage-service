package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/stretchr/testify/assert"
)

func initStorageAggregate(storages []domain.Storage) *domain.SchoolStorageAggregate {
	aggregate := domain.NewSchoolStorageAggregate()
	aggregate.Storages = storages
	return aggregate
}

func TestAddStorage(t *testing.T) {
	tests := []struct {
		name            string
		storages        []domain.Storage
		storageName     string
		storageLocation string
		err             error
		expectError     bool
	}{
		{
			name:            "add storage",
			storages:        []domain.Storage{},
			storageName:     "storage",
			storageLocation: "location",
			err:             nil,
			expectError:     false,
		},
		{
			name:            "add storage without name",
			storages:        []domain.Storage{},
			storageName:     "",
			storageLocation: "location",
			err:             domain.ErrStorageNameNotSet,
			expectError:     true,
		},
		{
			name:            "add storage without location",
			storages:        []domain.Storage{},
			storageName:     "storage",
			storageLocation: "",
			err:             domain.ErrStorageLocationNotSet,
			expectError:     true,
		},
		{
			name: "storage already exists",
			storages: []domain.Storage{{
				ID:       uuid.NewString(),
				Name:     "storage",
				Location: "location",
			}},
			storageName:     "storage",
			storageLocation: "location",
			err:             domain.ErrStorageAlreadyExists("storage", "location"),
			expectError:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initStorageAggregate(test.storages)
			storageID, err := aggregate.AddStorage(test.storageName, test.storageLocation)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.NotEqual(t, "", storageID)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, domain.StorageAdded, event.EventType())
		})
	}
}

func TestRemoveStorage(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name        string
		storages    []domain.Storage
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "remove storage",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name:        "error when removing not existing storage",
			storages:    []domain.Storage{},
			reason:      "test",
			err:         domain.ErrStorageIDNotFound(storageID),
			expectError: true,
		},
		{
			name: "error when removing without a reason",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initStorageAggregate(test.storages)
			err := aggregate.RemoveStorage(storageID, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.Len(t, aggregate.Storages, 0)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, domain.StorageRemoved, event.EventType())
		})
	}
}

func TestRenameStorage(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name        string
		storages    []domain.Storage
		storageName string
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "rename storage",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "renamed",
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name: "storage with same name and location exists",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "storage",
			reason:      "test",
			err:         domain.ErrStorageAlreadyExists("storage", "location"),
			expectError: true,
		},
		{
			name: "storage name not set",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "",
			reason:      "test",
			err:         domain.ErrStorageNameNotSet,
			expectError: true,
		},
		{
			name: "reason not specified",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "renamed",
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initStorageAggregate(test.storages)
			err := aggregate.RenameStorage(storageID, test.storageName, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, domain.StorageRenamed, event.EventType())
		})
	}
}

func TestSetStorageLocation(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name            string
		storages        []domain.Storage
		storageLocation string
		reason          string
		err             error
		expectError     bool
	}{
		{
			name: "rename storage",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "relocated",
			reason:          "test",
			err:             nil,
			expectError:     false,
		},
		{
			name: "storage with same name and location exists",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "location",
			reason:          "test",
			err:             domain.ErrStorageAlreadyExists("storage", "location"),
			expectError:     true,
		},
		{
			name: "storage name not set",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "",
			reason:          "test",
			err:             domain.ErrStorageLocationNotSet,
			expectError:     true,
		},
		{
			name: "reason not specified",
			storages: []domain.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "relocated",
			reason:          "",
			err:             domain.ErrReasonNotSpecified,
			expectError:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initStorageAggregate(test.storages)
			err := aggregate.RelocateStorage(storageID, test.storageLocation, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, domain.StorageRelocated, event.EventType())
		})
	}
}
