package main

import (
	"context"
	"log"
	"net/http"

	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/kammeph/school-book-storage-service/web"
)

func main() {
	// postgres.CreatePostgresStoreTables()
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
	// web.InMemoryConfig()
	web.PostgresMongoRabbitConfig(db, client, connection)
	// http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	web.JsonResponse(w, "hello")
	// })
	http.ListenAndServe(":9090", nil)
}
