package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/web"
)

type AuthController struct {
	commands userapp.UserCommandHandlers
	queries  userapp.UserQueryHandlers
}

func NewAuthController(commands userapp.UserCommandHandlers, queries userapp.UserQueryHandlers) *AuthController {
	return &AuthController{commands, queries}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var command userapp.LoginUserCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	user, err := c.commands.LoginUserHandler.Handle(context.Background(), command)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	accessToken, err := createAccessToken(*user, jwtSecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	refreshToken, err := createRefreshToken(user.ID, jwtSecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	web.JsonResponse(w, accessToken)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var command userapp.LoginUserCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := c.commands.RegisterUserHandler.Handle(context.Background(), userapp.RegisterUserCommand(command)); err != nil {
		http.Error(w, err.Error(), 400)
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getRefreshToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	claims := &RefreshClaims{}
	if err := getClaimsFromToken(r, tokenString, claims); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := c.queries.GetUserByIDHandler.Handle(
		context.Background(),
		userapp.NewGetUserByIDQuery("users", claims.UserID),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accessToken, err := createAccessToken(*user, jwtSecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	web.JsonResponse(w, accessToken)
}
