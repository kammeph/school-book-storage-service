package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
)

var (
	jwtSecretKey             = utils.GetenvOrFallback("JWT_SECRET_KEY", "MySuperSecretKey")
	jwtAccessTokenExpiry, _  = strconv.Atoi(utils.GetenvOrFallback("JWT_ACCESS_TOKEN_EXPIRY_SEC", "60"))
	jwtRefreshTokenExpiry, _ = strconv.Atoi(utils.GetenvOrFallback("JWT_REFRESH_TOKEN_EXPIRY_SEC", "120"))
)

type AccessClaims struct {
	jwt.StandardClaims
	UserID   string            `json:"userId"`
	UserName string            `json:"userName"`
	Roles    []userdomain.Role `json:"roles"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}

func IsAllowed(handler func(w http.ResponseWriter, r *http.Request), roles []userdomain.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getAccessToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims := &AccessClaims{}
		if err := getClaimsFromToken(r, tokenString, claims); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		for _, role := range roles {
			if fp.Some(claims.Roles, func(r userdomain.Role) bool { return r == role }) {
				handler(w, r)
				return
			}
		}
		http.Error(w, "user missing permissions", http.StatusMethodNotAllowed)
	}
}

func IsAllowedWithClaims(
	handler func(w http.ResponseWriter, r *http.Request, claims AccessClaims),
	roles []userdomain.Role,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getAccessToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims := &AccessClaims{}
		if err := getClaimsFromToken(r, tokenString, claims); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		for _, role := range roles {
			if fp.Some(claims.Roles, func(r userdomain.Role) bool { return r == role }) {
				handler(w, r, *claims)
				return
			}
		}
		http.Error(w, "user missing permissions", http.StatusForbidden)
	}
}

func getAccessToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("access token is not set")
	}
	if !strings.ContainsAny(auth, "Bearer") {
		return "", errors.New("no bearer token found")
	}
	token := strings.Split(auth, " ")[1]
	return token, nil
}

func getRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return "", err
	}
	// if cookie.Valid() != nil {
	// 	return "", cookie.Valid()
	// }
	return cookie.Value, nil
}

func getClaimsFromToken(r *http.Request, tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("access token is invalid")
	}
	return nil
}

func createAccessToken(user userdomain.UserModel, secret string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(jwtAccessTokenExpiry) * time.Second)
	claims := AccessClaims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    user.Name,
		},
		user.ID, user.Name, user.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func createRefreshToken(userID string, secret string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(jwtRefreshTokenExpiry) * time.Second)
	claims := RefreshClaims{
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
