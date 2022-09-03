package storages

import (
	"database/sql"

	"github.com/kammeph/school-book-storage-service/application/storageapp"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/kammeph/school-book-storage-service/web"
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
	web.Get(
		"/api/storages/get-all/",
		web.IsAllowed(
			controller.GetAllStorages,
			[]userdomain.Role{userdomain.User, userdomain.Superuser, userdomain.Admin},
		))
	web.Get(
		"/api/storages/get-by-id/",
		web.IsAllowed(
			controller.GetStorageByID,
			[]userdomain.Role{userdomain.User, userdomain.Superuser, userdomain.Admin},
		))
	web.Get(
		"/api/storages/get-by-name/",
		web.IsAllowed(
			controller.GetStorageByName,
			[]userdomain.Role{userdomain.User, userdomain.Superuser, userdomain.Admin},
		))
	web.Post(
		"/api/storages/add",
		web.IsAllowed(
			controller.AddStorage,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/storages/remove",
		web.IsAllowed(
			controller.RemoveStorage,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/storages/rename",
		web.IsAllowed(
			controller.RenameStorage,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/storages/relocate",
		web.IsAllowed(
			controller.RelocateStorage,
			[]userdomain.Role{userdomain.Admin},
		))
}
