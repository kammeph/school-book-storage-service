package main

import (
	"context"
	"log"
	"net/http"

	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/kammeph/school-book-storage-service/web/auth"
	"github.com/kammeph/school-book-storage-service/web/school"
	"github.com/kammeph/school-book-storage-service/web/storages"
	"github.com/kammeph/school-book-storage-service/web/users"
)

func main() {
	connection := rabbitmq.NewRabbitMQConnection()
	defer func() {
		if err := connection.Close(); err != nil {
			panic(err)
		}
		log.Println("Connection to rabbit mq closed.")
	}()
	db := postgresdb.NewPostgresDB()
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
		log.Println("Connection to postgres db closed.")
	}()
	client := mongodb.NewMongoClient()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		log.Println("Connection to mongo db closed.")
	}()
	auth.PostgresConfig(db)
	users.PostgresConfig(db)
	school.PostgresMongoRabbitConfig(db, client, connection)
	storages.PostgresMongoRabbitConfig(db, client, connection)
	http.ListenAndServe(":9090", nil)
}
