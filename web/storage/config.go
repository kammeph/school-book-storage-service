package storage

import (
	"database/sql"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgres"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	infrastructure "github.com/kammeph/school-book-storage-service/infrastructure/storage"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigureEndpointWithPostgresStore(
	postgresdb *sql.DB,
	mongodb *mongo.Client,
	connection rabbitmq.AmqpConnection,
) {
	publisher, err := rabbitmq.NewRabbitEventPublisher(connection, "storage")
	if err != nil {
		panic(err)
	}
	subscriber, err := rabbitmq.NewRabbitEventSubscriber(connection)
	if err != nil {
		panic(err)
	}
	store := postgres.NewPostgressStore("storages", postgresdb)
	repository := infrastructure.NewStorageWithBookRepository(mongodb, "school_book_storage", "storages")

	eventHandler := storage.NewStorageEventHandler(repository)
	subscriber.Subscribe("storage", eventHandler)
	subscriber.Subscribe("storage", &storage.TestHandler{})

	commandHandlers := storage.NewStorageCommandHandlers(store, publisher)
	queryHandlers := storage.NewStorageQueryHandlers(repository)

	controller := NewStorageController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func configureEndpoints(controller *StorageController) {
	http.HandleFunc("/api/storages/get-all/", controller.GetAllStorages)
	http.HandleFunc("/api/storages/get-by-id/", controller.GetStorageByID)
	http.HandleFunc("/api/storages/get-by-name/", controller.GetStorageByName)
	http.HandleFunc("/api/storages/add", controller.AddStorage)
	http.HandleFunc("/api/storages/remove", controller.RemoveStorage)
	http.HandleFunc("/api/storages/rename", controller.SetStorageName)
	http.HandleFunc("/api/storages/relocate", controller.SetStorageLocation)
}
