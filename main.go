package main

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/infrastructure/dbs"
	"github.com/kammeph/school-book-storage-service/infrastructure/messagebroker"
	"github.com/kammeph/school-book-storage-service/web/storage"
)

func main() {
	connection, err := messagebroker.NewRabbitMQConnection()
	defer connection.Close()
	if err != nil {
		panic(err)
	}
	db := dbs.NewPostgresDB()
	defer db.Close()
	storage.ConfigureEndpointWithPostgresStore(db, connection)
	// storage.ConfigureEndpointsWithMemoryStore()
	http.ListenAndServe(":9090", nil)
}
