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
		ID       string            `json:"id"`
		SchoolID string            `json:"schoolId"`
		Name     string            `json:"name"`
		Roles    []userdomain.Role `json:"roles"`
		Locale   userdomain.Locale `json:"locale"`
	}{
		ID:       claims.UserID,
		SchoolID: claims.SchoolID,
		Name:     claims.UserName,
		Roles:    claims.Roles,
		Locale:   claims.Locale,
	}
	web.JsonResponse(w, user)
}
