package storages

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kammeph/school-book-storage-service/application/storageapp"
	"github.com/kammeph/school-book-storage-service/web"
)

type StorageController struct {
	commmandHandlers storageapp.StorageCommandHandlers
	queryHandlers    storageapp.StorageQueryHandlers
}

func NewStorageController(commandHandlers storageapp.StorageCommandHandlers, queryHandlers storageapp.StorageQueryHandlers) *StorageController {
	return &StorageController{commandHandlers, queryHandlers}
}

func (c StorageController) AddStorage(w http.ResponseWriter, r *http.Request) {
	var command storageapp.AddStorageCommand
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	storageID, err := c.commmandHandlers.AddStorageHandler.Handle(ctx, command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	web.HttpResponse(w, storageID)
}

func (c StorageController) RemoveStorage(w http.ResponseWriter, r *http.Request) {
	var command storageapp.RemoveStorageCommand
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.RemoveStorageHandler.Handle(ctx, command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c StorageController) RenameStorage(w http.ResponseWriter, r *http.Request) {
	var command storageapp.RenameStorageCommand
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.RenameStorageHandler.Handle(ctx, command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c StorageController) RelocateStorage(w http.ResponseWriter, r *http.Request) {
	var command storageapp.RelocateStorageCommand
	json.NewDecoder(r.Body).Decode(&command)
	ctx := context.Background()
	defer ctx.Done()
	err := c.commmandHandlers.RelocateStorageHandler.Handle(ctx, command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c StorageController) GetAllStorages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-1]
	query := storageapp.NewGetAllStorages(aggregateID)
	storages, err := c.queryHandlers.GetAllHandler.Handle(ctx, query)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	web.HttpResponse(w, storages)
}

func (c StorageController) GetStorageByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-2]
	storageID := path[len(path)-1]
	query := storageapp.NewGetStorageByID(aggregateID, storageID)
	storage, err := c.queryHandlers.GetStorageByIDHandler.Handle(ctx, query)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	web.HttpResponse(w, storage)
}

func (c StorageController) GetStorageByName(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	aggregateID := path[len(path)-2]
	name := path[len(path)-1]
	query := storageapp.NewGetStorageByName(aggregateID, name)
	storage, err := c.queryHandlers.GetStorageByNameHandler.Handle(ctx, query)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	web.HttpResponse(w, storage)
}
