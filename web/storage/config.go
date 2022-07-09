package storage

import (
	"database/sql"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/kammeph/school-book-storage-service/infrastructure/stores"
)

func ConfigureEndpointsWithMemoryStore() {
	store := stores.NewMemoryStore()
	serializer := serializers.NewJSONSerializer()
	repository := storage.NewStorageRepository(store, serializer)
	configureEndpoints(repository)
}

func ConfigureEndpointWithPostgresStore(db *sql.DB) {
	store := stores.NewPostgressStore("storages", db)
	serializer := serializers.NewJSONSerializer()
	repository := storage.NewStorageRepository(store, serializer)
	configureEndpoints(repository)
}

func configureEndpoints(repository *common.Repository) {
	controller := NewStorageController(repository)
	http.HandleFunc("/api/storages/get-all/", controller.GetAllStorages)
	http.HandleFunc("/api/storages/get-by-id/", controller.GetStorageByID)
	http.HandleFunc("/api/storages/get-by-name/", controller.GetStorageByName)
	http.HandleFunc("/api/storages/add", controller.AddStorage)
	http.HandleFunc("/api/storages/remove", controller.RemoveStorage)
	http.HandleFunc("/api/storages/set-name", controller.SetStorageName)
	http.HandleFunc("/api/storages/set-location", controller.SetStorageLocation)
}
