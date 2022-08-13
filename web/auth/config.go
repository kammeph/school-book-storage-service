package auth

import (
	"database/sql"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/web"
)

func PostgresConfig(db *sql.DB) {
	store := postgresdb.NewPostgresStore("users", db)
	commandHandlers := userapp.NewUsersCommandHandlers(store, nil)
	queryHandlers := userapp.NewUserQueryHandlers(store)
	controller := NewAuthController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func configureEndpoints(controller *AuthController) {
	web.Post("/api/auth/login", controller.Login)
	web.Post("/api/auth/register", controller.Register)
	web.Get("/api/auth/refresh", controller.Refresh)
}
