package main

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/infrastructure/dbs"
	"github.com/kammeph/school-book-storage-service/web/storage"
)

func main() {
	db := dbs.NewPostgresDB()
	defer db.Close()
	storage.ConfigureEndpointWithPostgresStore(db)
	// storage.ConfigureEndpointsWithMemoryStore()
	http.ListenAndServe(":9090", nil)
}
