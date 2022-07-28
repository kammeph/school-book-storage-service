package postgresdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	_ "github.com/lib/pq"
)

var (
	pgdriver   = utils.GetenvOrFallback("PG_DRIVER", "postgres")
	pguser     = utils.GetenvOrFallback("PG_USER", "postgres")
	pgpassword = utils.GetenvOrFallback("PG_PASSWORD", "rootpwd")
	pghost     = utils.GetenvOrFallback("PG_HOST", "localhost")
	pgport     = utils.GetenvOrFallback("PG_PORT", "5432")
	pgdbname   = utils.GetenvOrFallback("PG_DATABASE", "school_book_storage")
	pgsslmode  = utils.GetenvOrFallback("PG_SSLMODE", "disable")
)

func NewPostgresDB() *sql.DB {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		pguser, pgpassword, pghost, pgport, pgdbname, pgsslmode)
	db, err := sql.Open(pgdriver, connStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	log.Println("Successfully connected and pinged to postgres db.")
	return db
}
