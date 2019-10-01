package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type UserAuth struct {
	Username  string
	Roles     []string
	IPAddress string
	Token     string
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}
type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromHttpRequest(r)

			username, _ := UsernameFromToken(token)

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			userAuth := UserAuth{
				Username:  username,
				IPAddress: ip,
				Token:     token,
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &userAuth)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
func TokenFromHttpRequest(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}

func UsernameFromToken(tokenString string) (string, error) {
	token, err := JwtDecode(tokenString)

	if err != nil {
		return "", errors.New("access_denied")
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims == nil {
			return "", errors.New("access_denied")
		}
		return claims.Username, nil
	} else {
		return "", errors.New("access_denied")
	}
}

func ForContext(ctx context.Context) *UserAuth {
	raw := ctx.Value(userCtxKey)
	if raw == nil {
		return nil
	}
	return raw.(*UserAuth)
}

func GetAuthFromContext(ctx context.Context) *UserAuth {
	return ForContext(ctx)
}
