package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mySigningKey = []byte("table.az")

func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
}

func JwtCreate(userID primitive.ObjectID, expiresAt int64) string {
	claims := UserClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "table",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)

	return ss
}
