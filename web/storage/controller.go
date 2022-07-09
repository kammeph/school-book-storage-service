package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	application "github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
	"github.com/kammeph/school-book-storage-service/web/common"
)

type StorageController struct {
	commmandHandlers storage.StorageCommandHandlers
	queryHandlers    storage.StorageQueryHandlers
}

func NewStorageController(repository *application.Repository) *StorageController {
	return &StorageController{
		commmandHandlers: storage.NewStorageCommandHandlers(repository),
		queryHandlers:    storage.NewStorageQueryHandlers(repository),
	}
}

func (c StorageController) AddStorage(w http.ResponseWriter, r *http.Request) {
	var command storage.AddStorage
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	dto, err := c.commmandHandlers.AddStorageHandler.Handle(ctx, command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	common.JsonResponse(w, dto)
}

func (c StorageController) RemoveStorage(w http.ResponseWriter, r *http.Request) {
	var command storage.RemoveStorage
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.RemoveStorageHandler.Handle(ctx, command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (c StorageController) SetStorageName(w http.ResponseWriter, r *http.Request) {
	var command storage.SetStorageName
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.SetStorageNameHandler.Handle(ctx, command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (c StorageController) SetStorageLocation(w http.ResponseWriter, r *http.Request) {
	var command storage.SetStorageLocation
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.SetStorageLocationHandler.Handle(ctx, command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (c StorageController) GetAllStorages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-1]
	query := storage.NewGetAllStorages(aggregateID)
	storages, err := c.queryHandlers.GetAllHandler.Handle(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	common.JsonResponse(w, storages)
}

func (c StorageController) GetStorageByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-2]
	storageID := path[len(path)-1]
	query := storage.NewGetStorageByID(aggregateID, storageID)
	storage, err := c.queryHandlers.GetStorageByIDHandler.Handle(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	common.JsonResponse(w, storage)
}

func (c StorageController) GetStorageByName(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-2]
	name := path[len(path)-1]
	query := storage.NewGetStorageByName(aggregateID, name)
	storage, err := c.queryHandlers.GetStorageByNameHandler.Handle(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	common.JsonResponse(w, storage)
}
