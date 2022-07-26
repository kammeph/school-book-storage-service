package main

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/web"
)

func main() {
	// postgres.CreatePostgresStoreTables()
	//web.InMemoryConfig()
	web.PostgresMongoRabbitConfig()
	http.ListenAndServe(":9090", nil)
}
