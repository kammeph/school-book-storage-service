package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kammeph/school-book-storage-service/application/userapp"
	"github.com/kammeph/school-book-storage-service/web"
)

type AccessTokenResponseModel struct {
	AccessToken string `json:"accessToken"`
}

func AccessTokenResponse(w http.ResponseWriter, accessToken string) {
	web.JsonResponse(w, AccessTokenResponseModel{accessToken})
}

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
		web.HttpErrorResponse(w, err.Error())
		return
	}
	user, err := c.commands.LoginUserHandler.Handle(context.Background(), command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	accessToken, err := web.CreateAccessToken(*user)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	refreshToken, err := web.CreateRefreshToken(user.ID)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	AccessTokenResponse(w, accessToken)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var command userapp.RegisterUserCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.commands.RegisterUserHandler.Handle(context.Background(), userapp.RegisterUserCommand(command)); err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := web.GetRefreshToken(r)
	if err != nil {
		web.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	claims := &web.RefreshClaims{}
	if err := web.GetClaimsFromToken(r, tokenString, claims); err != nil {
		web.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := c.queries.GetUserByIDHandler.Handle(
		context.Background(),
		userapp.NewGetUserByIDQuery("users", claims.UserID),
	)
	if err != nil {
		web.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accessToken, err := web.CreateAccessToken(*user)
	if err != nil {
		web.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	AccessTokenResponse(w, accessToken)
}
