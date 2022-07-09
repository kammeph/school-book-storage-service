package storage_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	application "github.com/kammeph/school-book-storage-service/application/storage"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/kammeph/school-book-storage-service/infrastructure/stores"
	"github.com/kammeph/school-book-storage-service/web/storage"
	"github.com/stretchr/testify/assert"
)

func getStorageController() (storage.StorageController, string, error) {
	store := stores.NewMemoryStore()
	serializer := serializers.NewJSONSerializer()
	repository := application.NewStorageRepository(store, serializer)
	aggregate, err := repository.Load(context.Background(), "test")
	if err != nil {
		return storage.StorageController{}, "", err
	}
	schoolAggregate, _ := aggregate.(*domain.SchoolAggregateRoot)
	storageID, err := schoolAggregate.AddStorage("storage", "location")
	if err != nil {
		return storage.StorageController{}, "", err
	}
	repository.Save(context.Background(), schoolAggregate)
	controller := storage.NewStorageController(repository)
	return *controller, storageID, nil
}

func TestAddStorage(t *testing.T) {
	controller, _, _ := getStorageController()
	body := `{"aggregateId":"storage","name":"storage","location":"location"}`
	req, err := http.NewRequest("POST", "api/storages/add", strings.NewReader(body))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.AddStorage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Contains(t, rr.Body.String(), "storageId")
}

func TestRemoveStorage(t *testing.T) {
	controller, storageID, err := getStorageController()
	body := fmt.Sprintf(`{"aggregateId":"test","storageId":"%s","reason":"test"}`, storageID)
	req, err := http.NewRequest("POST", "api/storages/remove", strings.NewReader(body))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.RemoveStorage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestSetStorageName(t *testing.T) {
	controller, storageID, err := getStorageController()
	body := fmt.Sprintf(`{"aggregateId":"test","storageId":"%s","name":"storage set name","reason":"test"}`, storageID)
	req, err := http.NewRequest("POST", "api/storages/set-name", strings.NewReader(body))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.SetStorageName)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestSetStorageLocation(t *testing.T) {
	controller, storageID, err := getStorageController()
	body := fmt.Sprintf(`{"aggregateId":"test","storageId":"%s","location":"location set","reason":"test"}`, storageID)
	req, err := http.NewRequest("POST", "api/storages/set-location", strings.NewReader(body))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.SetStorageLocation)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestGetAllStorages(t *testing.T) {
	controller, storageID, err := getStorageController()
	req, err := http.NewRequest("GET", "api/storages/get-all/test", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetAllStorages)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"ID":"%s"`, storageID))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Name":"%s"`, "storage"))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Location":"%s"`, "location"))
}

func TestGetStorageById(t *testing.T) {
	controller, storageID, err := getStorageController()
	req, err := http.NewRequest("GET", fmt.Sprintf("api/storages/get-by-id/test/%s", storageID), nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetStorageByID)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"ID":"%s"`, storageID))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Name":"%s"`, "storage"))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Location":"%s"`, "location"))
}

func TestGetStorageByName(t *testing.T) {
	controller, storageID, err := getStorageController()
	req, err := http.NewRequest("GET", fmt.Sprintf("api/storages/get-by-name/test/%s", "storage"), nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetStorageByName)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"ID":"%s"`, storageID))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Name":"%s"`, "storage"))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"Location":"%s"`, "location"))
}
