package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/storage"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongo"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgres"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	infrastructure "github.com/kammeph/school-book-storage-service/infrastructure/storage"
)

func InMemoryConfig() {
	broker := memory.NewMemoryMessageBroker()
	store := memory.NewMemoryStore()
	repository := infrastructure.NewMemoryRepository()

	eventHandler := storage.NewStorageEventHandler(repository)
	broker.Subscribe("storage", eventHandler)
	broker.Subscribe("storage", &storage.TestHandler{})

	commandHandlers := storage.NewStorageCommandHandlers(store, broker)
	queryHandlers := storage.NewStorageQueryHandlers(repository)

	controller := NewStorageController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func PostgresMongoRabbitConfig() {
	connection := rabbitmq.NewRabbitMQConnection()
	defer func() {
		if err := connection.Close(); err != nil {
			panic(err)
		}
		fmt.Println("Connection to rabbit mq closed.")
	}()
	postgresdb := postgres.NewDB()
	defer func() {
		if err := postgresdb.Close(); err != nil {
			panic(err)
		}
		fmt.Println("Connection to postgres db closed.")
	}()
	mongodb := mongo.NewDB()
	defer func() {
		if err := mongodb.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("Connection to mongo db closed.")
	}()
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
	http.HandleFunc("/api/storages/rename", controller.RenameStorage)
	http.HandleFunc("/api/storages/relocate", controller.RelocateStorage)
}
