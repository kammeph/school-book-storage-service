package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kammeph/school-book-storage-service/infrastructure/mongo"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgres"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/kammeph/school-book-storage-service/web/storage"
)

func main() {
	// postgres.CreatePostgresStoreTables()
	rabbitMqConnection := rabbitmq.NewRabbitMQConnection()
	defer func() {
		if err := rabbitMqConnection.Close(); err != nil {
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
	storage.ConfigureEndpointWithPostgresStore(postgresdb, mongodb, rabbitMqConnection)
	// storage.ConfigureEndpointsWithMemoryStore()
	http.ListenAndServe(":9090", nil)
}
