package main

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/web/storage"
)

func main() {
	storage.ConfigureEndpointsWithMemoryStore()
	http.ListenAndServe(":9090", nil)
}
