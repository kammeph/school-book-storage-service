package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
)

func InMemoryConfig() {
	broker := memory.NewMemoryMessageBroker()
	store := memory.NewMemoryStore()
	repository := memory.NewMemoryRepository()

	eventHandler := application.NewStorageEventHandler(repository)
	broker.Subscribe("storage", eventHandler)
	broker.Subscribe("storage", &application.TestHandler{})

	commandHandlers := application.NewStorageCommandHandlers(store, broker)
	queryHandlers := application.NewStorageQueryHandlers(repository)

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
	db := postgresdb.NewPostgresDB()
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
		fmt.Println("Connection to postgres db closed.")
	}()
	client := mongodb.NewMongoClient()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
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
	store := postgresdb.NewPostgresStore("storages", db)
	repository := mongodb.NewStorageWithBookRepository(client, "school_book_storage", "storages")

	eventHandler := application.NewStorageEventHandler(repository)
	subscriber.Subscribe("storage", eventHandler)
	subscriber.Subscribe("storage", &application.TestHandler{})

	commandHandlers := application.NewStorageCommandHandlers(store, publisher)
	queryHandlers := application.NewStorageQueryHandlers(repository)

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
