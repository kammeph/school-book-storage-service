package school

import (
	"database/sql"

	"github.com/kammeph/school-book-storage-service/application/schoolapp"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/kammeph/school-book-storage-service/web"
)

func PostgresMongoRabbitConfig(postgresDB *sql.DB, mongoClient mongodb.Client, rabbit rabbitmq.AmqpConnection) {
	publisher, err := rabbitmq.NewRabbitEventPublisher(rabbit, "school")
	if err != nil {
		panic(err)
	}
	subscriber, err := rabbitmq.NewRabbitEventSubscriber(rabbit)
	if err != nil {
		panic(err)
	}

	store := postgresdb.NewPostgresStore("schools", postgresDB)
	repository := mongodb.NewSchoolRepository(mongoClient, "school_book_storage", "schools")

	eventHandler := schoolapp.NewSchoolEventHandler(repository)
	subscriber.Subscribe("school", eventHandler)

	commandHandlers := schoolapp.NewSchoolCommandHandlers(store, publisher)
	queryHandlers := schoolapp.NewSchoolQueryHandlers(repository)

	controller := NewSchoolController(*commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func configureEndpoints(controller *SchoolController) {
	web.Get(
		"/api/schools/get-all",
		web.IsAllowed(
			controller.GetSchools,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Get(
		"/api/schools/get-by-id/",
		web.IsAllowed(
			controller.GetSchoolByID,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/schools/add",
		web.IsAllowed(
			controller.AddSchool,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/schools/deactivate",
		web.IsAllowed(
			controller.DeactivateSchool,
			[]userdomain.Role{userdomain.Admin},
		))
	web.Post(
		"/api/schools/rename",
		web.IsAllowed(
			controller.RenameSchool,
			[]userdomain.Role{userdomain.Admin},
		))
}
