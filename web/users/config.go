package users

import (
	"database/sql"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/kammeph/school-book-storage-service/web"
)

func PostgresConfig(db *sql.DB) {
	store := postgresdb.NewPostgresStore("users", db)
	commandHandlers := userapp.NewUsersCommandHandlers(store, nil)
	queryHandlers := userapp.NewUserQueryHandlers(store)
	controller := NewUsersController(commandHandlers, queryHandlers)
	configureEndpoints(controller)
}

func configureEndpoints(controller *UsersController) {
	web.Get(
		"/api/users/me",
		web.IsAllowedWithClaims(
			controller.GetMe,
			[]userdomain.Role{userdomain.User, userdomain.Superuser, userdomain.Admin},
		))
}
