package helpers

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/kamilniftaliev/table-server/api/models"
)

var ErrorAccessDenied = errors.New("access_denied")

type contextKey struct {
	name string
}

var authCtxKey = &contextKey{"auth"}

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func getAuthFromRequest(r *http.Request) *models.Auth {
	token := getTokenFromRequest(r)

	username := getUsernameFromToken(token)
	var authError error = nil

	if username == "" {
		authError = ErrorAccessDenied
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	auth := models.Auth{
		Username:  username,
		IPAddress: ip,
		Token:     token,
		Error:     authError,
	}

	return &auth
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := getAuthFromRequest(r)

			// put it in context
			ctx := context.WithValue(r.Context(), authCtxKey, auth)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenFromRequest(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}

func getUsernameFromToken(tokenString string) string {
	token, err := JwtDecode(tokenString)

	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims == nil {
			return ""
		}
		return claims.Username
	} else {
		return ""
	}
}

func GetAuth(ctx context.Context) *models.Auth {
	return ctx.Value(authCtxKey).(*models.Auth)
}
