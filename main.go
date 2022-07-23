package main

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/web/storage"
)

func main() {
	// postgres.CreatePostgresStoreTables()
	storage.InMemoryConfig()
	// storage.PostgresMongoRabbitConfig()
	http.ListenAndServe(":9090", nil)
}
