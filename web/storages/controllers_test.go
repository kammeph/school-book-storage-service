package storages_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/application/storageapp"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/kammeph/school-book-storage-service/web/storages"
	"github.com/stretchr/testify/assert"
)

func getStorageController() *storages.StorageController {
	eventDataForRemove, _ := json.Marshal(storagedomain.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testRemove",
		Name:      "Closet to Remove",
		Location:  "Room 12",
	})
	eventForRemove := domain.EventModel{
		ID:      "school",
		Type:    storagedomain.StorageAdded,
		Version: 1,
		At:      time.Now(),
		Data:    string(eventDataForRemove),
	}
	eventDataForUpdate, _ := json.Marshal(storagedomain.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testUpdate",
		Name:      "Closet to Update",
		Location:  "Room 12",
	})
	eventForUpdate := domain.EventModel{
		ID:      "school",
		Type:    storagedomain.StorageAdded,
		Version: 2,
		At:      time.Now(),
		Data:    string(eventDataForUpdate),
	}
	store := memory.NewMemoryStoreWithEvents([]domain.Event{&eventForRemove, &eventForUpdate})
	storage1School1 := storagedomain.NewStorageWithBooks("school1", "storage1School1", "Closet 1", "Room 101")
	storage2School1 := storagedomain.NewStorageWithBooks("school1", "storage2School1", "Closet 2", "Room 101")
	storage1School2 := storagedomain.NewStorageWithBooks("school2", "storage1School2", "Closet 1", "Room 203")
	repository := memory.NewMemoryRepositoryWithStorages(
		[]storagedomain.StorageWithBooks{storage1School1, storage2School1, storage1School2})
	commandHandlers := storageapp.NewStorageCommandHandlers(store, nil)
	queryHandlers := storageapp.NewStorageQueryHandlers(repository)
	return storages.NewStorageController(commandHandlers, queryHandlers)
}

func TestAddStorage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "add storage success",
			body:       `{"aggregateId":"school","name":"storage","location":"location"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "add storage failed",
			body:       `{"aggregateId":"school","name":"Closet to Update","location":"Room 12"}`,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("POST", "api/storages/add", strings.NewReader(test.body))
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.AddStorage)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			assert.NotEqual(t, "", rr.Body.String())
		})
	}
}

func TestRemoveStorage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{
			name:       "remove storage success",
			body:       `{"aggregateId":"school","storageId":"testRemove","reason":"test"}`,
			statusCode: http.StatusOK,
			response:   "",
		},
		{
			name:       "remove storage failed",
			body:       `{"aggregateId":"school","storageId":"unknown","reason":"test"}`,
			statusCode: http.StatusBadRequest,
			response:   fmt.Sprintln(storagedomain.ErrStorageIDNotFound("unknown").Error()),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("POST", "api/storages/remove", strings.NewReader(test.body))
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.RemoveStorage)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.response, rr.Body.String())
			assert.Equal(t, test.statusCode, rr.Code)
		})
	}
}

func TestRenameStorage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{
			name:       "rename storage success",
			body:       `{"aggregateId":"school","storageId":"testUpdate","name":"storage set name","reason":"test"}`,
			statusCode: http.StatusOK,
			response:   "",
		},
		{
			name:       "rename storage failed",
			body:       `{"aggregateId":"school","storageId":"testUpdate","name":"Closet to Update","reason":"test"}`,
			statusCode: http.StatusBadRequest,
			response:   fmt.Sprintln(storagedomain.ErrStorageAlreadyExists("Closet to Update", "Room 12").Error()),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("POST", "api/storages/rename", strings.NewReader(test.body))
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.RenameStorage)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			assert.Equal(t, test.response, rr.Body.String())
		})
	}
}

func TestRelocateStorage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{
			name:       "rename storage success",
			body:       `{"aggregateId":"school","storageId":"testUpdate","location":"location set","reason":"test"}`,
			statusCode: http.StatusOK,
			response:   "",
		},
		{
			name:       "rename storage failed",
			body:       `{"aggregateId":"school","storageId":"testUpdate","location":"Room 12","reason":"test"}`,
			statusCode: http.StatusBadRequest,
			response:   fmt.Sprintln(storagedomain.ErrStorageAlreadyExists("Closet to Update", "Room 12").Error()),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("POST", "api/storages/relocate", strings.NewReader(test.body))
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.RelocateStorage)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			assert.Equal(t, test.response, rr.Body.String())
		})
	}
}

func TestGetAllStorages(t *testing.T) {
	tests := []struct {
		name        string
		aggregateID string
		statusCode  int
		responses   []string
	}{
		{
			name:        "get all storages empyty array",
			aggregateID: "school",
			statusCode:  http.StatusOK,
			responses:   []string{"[]"},
		},
		{
			name:        "get all storages for school 1",
			aggregateID: "school1",
			statusCode:  http.StatusOK,
			responses: []string{
				`"schoolId":"school1","storageId":"storage1School1","name":"Closet 1","location":"Room 101"`,
				`"schoolId":"school1","storageId":"storage2School1","name":"Closet 2","location":"Room 101"`,
			},
		},
		{
			name:        "get all storages for school 2",
			aggregateID: "school2",
			statusCode:  http.StatusOK,
			responses: []string{
				`"schoolId":"school2","storageId":"storage1School2","name":"Closet 1","location":"Room 203"`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("GET", fmt.Sprintf("api/storages/get-all/%s", test.aggregateID), nil)
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.GetAllStorages)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			for _, response := range test.responses {
				assert.Contains(t, rr.Body.String(), response)
			}
		})
	}
}

func TestGetStorageById(t *testing.T) {
	tests := []struct {
		name        string
		aggregateID string
		storageID   string
		statusCode  int
		response    string
	}{
		{
			name:        "get storage 1 school 1",
			aggregateID: "school1",
			storageID:   "storage1School1",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school1","storageId":"storage1School1","name":"Closet 1","location":"Room 101"`,
		},
		{
			name:        "get storage 2 school 1",
			aggregateID: "school1",
			storageID:   "storage2School1",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school1","storageId":"storage2School1","name":"Closet 2","location":"Room 101"`,
		},
		{
			name:        "get storage 1 school 2",
			aggregateID: "school2",
			storageID:   "storage1School2",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school2","storageId":"storage1School2","name":"Closet 1","location":"Room 203"`,
		},
		{
			name:        "get storage error",
			aggregateID: "school1",
			storageID:   "storage3School1",
			statusCode:  http.StatusBadRequest,
			response:    "no storage with ID storage3School1 found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("GET", fmt.Sprintf("api/storages/get-by-id/%s/%s", test.aggregateID, test.storageID), nil)
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.GetStorageByID)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			assert.Contains(t, rr.Body.String(), test.response)
		})
	}
}

func TestGetStorageByName(t *testing.T) {
	tests := []struct {
		name        string
		aggregateID string
		storageName string
		statusCode  int
		response    string
	}{
		{
			name:        "get storage 1 school 1",
			aggregateID: "school1",
			storageName: "Closet 1",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school1","storageId":"storage1School1","name":"Closet 1","location":"Room 101"`,
		},
		{
			name:        "get storage 2 school 1",
			aggregateID: "school1",
			storageName: "Closet 2",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school1","storageId":"storage2School1","name":"Closet 2","location":"Room 101"`,
		},
		{
			name:        "get storage 1 school 2",
			aggregateID: "school2",
			storageName: "Closet 1",
			statusCode:  http.StatusOK,
			response:    `"schoolId":"school2","storageId":"storage1School2","name":"Closet 1","location":"Room 203"`,
		},
		{
			name:        "get storage error",
			aggregateID: "school1",
			storageName: "storage3School1",
			statusCode:  http.StatusBadRequest,
			response:    "no storage with name storage3School1 found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := getStorageController()
			req, err := http.NewRequest("GET", fmt.Sprintf("api/storages/get-by-name/%s/%s", test.aggregateID, test.storageName), nil)
			assert.Nil(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.GetStorageByName)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
			assert.Contains(t, rr.Body.String(), test.response)
		})
	}
}
