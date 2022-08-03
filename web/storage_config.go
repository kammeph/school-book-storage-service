package web

import (
	"database/sql"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/storageapp"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
)

func InMemoryConfig() {
	broker := memory.NewMemoryMessageBroker()
	store := memory.NewMemoryStore()
	repository := memory.NewMemoryRepository()

	eventHandler := storageapp.NewStorageEventHandler(repository)
	broker.Subscribe("storage", eventHandler)
	broker.Subscribe("storage", &storageapp.TestHandler{})

	commandHandlers := storageapp.NewStorageCommandHandlers(store, broker)
	queryHandlers := storageapp.NewStorageQueryHandlers(repository)

	controller := NewStorageController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func PostgresMongoRabbitConfig(postgresDB *sql.DB, mongoClient mongodb.Client, rabbit rabbitmq.AmqpConnection) {
	publisher, err := rabbitmq.NewRabbitEventPublisher(rabbit, "storage")
	if err != nil {
		panic(err)
	}
	subscriber, err := rabbitmq.NewRabbitEventSubscriber(rabbit)
	if err != nil {
		panic(err)
	}
	store := postgresdb.NewPostgresStore("storages", postgresDB)
	repository := mongodb.NewStorageWithBookRepository(mongoClient, "school_book_storage", "storages")

	eventHandler := storageapp.NewStorageEventHandler(repository)
	subscriber.Subscribe("storage", eventHandler)
	subscriber.Subscribe("storage", &storageapp.TestHandler{})

	commandHandlers := storageapp.NewStorageCommandHandlers(store, publisher)
	queryHandlers := storageapp.NewStorageQueryHandlers(repository)

	controller := NewStorageController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func configureEndpoints(controller *StorageController) {
	http.HandleFunc("/api/storages/get-all/", controller.GetAllStorages)
	http.HandleFunc("/api/storages/get-by-id/", controller.GetStorageByID)
	http.HandleFunc("/api/storages/get-by-name/", controller.GetStorageByName)
	http.HandleFunc("/api/storages/add", controller.AddStorage)
	http.HandleFunc("/api/storages/remove", controller.RemoveStorage)
	http.HandleFunc("/api/storages/rename", controller.RenameStorage)
	http.HandleFunc("/api/storages/relocate", controller.RelocateStorage)
}
