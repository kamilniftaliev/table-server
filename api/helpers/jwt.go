package helpers

import "github.com/dgrijalva/jwt-go"

var mySigningKey = []byte("table.az")

func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
}

func JwtCreate(userID string, expiresAt int64) string {
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
