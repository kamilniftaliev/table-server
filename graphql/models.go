// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package table

type Table struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int    `json:"expiresAt"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Table    *Table `json:"table"`
}
