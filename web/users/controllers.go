package users

import (
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/web"
)

type UserResponseModel struct {
	User userapp.UserDto `json:"user"`
}

func UserResponse(w http.ResponseWriter, user userapp.UserDto) {
	response := UserResponseModel{user}
	web.JsonResponse(w, response)
}

type UsersController struct {
	commandHandlers userapp.UserCommandHandlers
	queryHandlers   userapp.UserQueryHandlers
}

func NewUsersController(commandHandlers userapp.UserCommandHandlers, queryHandlers userapp.UserQueryHandlers) *UsersController {
	return &UsersController{commandHandlers, queryHandlers}
}

func (c *UsersController) GetMe(w http.ResponseWriter, r *http.Request, claims web.AccessClaims) {
	user := userapp.UserDto{
		ID:       claims.UserID,
		SchoolID: claims.SchoolID,
		Name:     claims.UserName,
		Roles:    claims.Roles,
		Locale:   claims.Locale,
	}
	UserResponse(w, user)
}
