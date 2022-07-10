package dbs

import (
	"database/sql"
	"fmt"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	_ "github.com/lib/pq"
)

const (
	dbDriver   = "DB_DRIVER"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbName     = "DB_NAME"
	dbSslmode  = "DB_SSLMODE"
)

var (
	driver   = utils.GetenvOrFallback(dbDriver, "postgres")
	user     = utils.GetenvOrFallback(dbUser, "postgres")
	password = utils.GetenvOrFallback(dbPassword, "rootpwd")
	host     = utils.GetenvOrFallback(dbHost, "localhost")
	port     = utils.GetenvOrFallback(dbPort, "5432")
	dbname   = utils.GetenvOrFallback(dbName, "school_book_storage")
	sslmode  = utils.GetenvOrFallback(dbSslmode, "disable")
)

func NewPostgresDB() *sql.DB {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode)
	db, err := sql.Open(driver, connStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
