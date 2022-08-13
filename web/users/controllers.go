package users

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/web"
	"github.com/kammeph/school-book-storage-service/web/auth"
)

type UsersController struct {
	commandHandlers userapp.UserCommandHandlers
	queryHandlers   userapp.UserQueryHandlers
}

func NewUsersController(commandHandlers userapp.UserCommandHandlers, queryHandlers userapp.UserQueryHandlers) *UsersController {
	return &UsersController{commandHandlers, queryHandlers}
}

func (c *UsersController) GetMe(w http.ResponseWriter, r *http.Request, claims auth.AccessClaims) {
	user := struct {
		ID    string            `json:"id"`
		Name  string            `json:"name"`
		Roles []userdomain.Role `json:"roles"`
	}{
		ID:    claims.UserID,
		Name:  claims.UserName,
		Roles: claims.Roles,
	}
	web.JsonResponse(w, user)
}
