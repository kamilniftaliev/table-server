package api

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*Token, error) {
	userAuth := GetAuthFromContext(ctx)

	if len(userAuth.Username) > 1 {
		token := Token{
			Token: userAuth.Token,
		}
		return &token, nil
	}

	user, err := FindUsername(ctx, username)

	if err != nil || user == nil || !ComparePassword(password, user.Password) {
		return nil, errors.New("error_password")
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()

	token := &Token{
		Token:     JwtCreate(user.Username, expiresAt),
		ExpiresAt: int(expiresAt),
	}

	return token, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) (*User, error) {
	auth := GetAuthFromContext(ctx)
	username, authErr := UsernameFromToken(auth.Token)

	if authErr != nil {
		return nil, authErr
	}

	usersCollection := DB.Collection("users")

	cur, curErr := usersCollection.Find(
		ctx,
		bson.M{"username": username},
	)

	log.Println("cur::::", cur)
	log.Println("curErr::::", curErr)

	// var user *User

	// err := cur.Decode(&user)

	// jsonErr := json.Unmarshal([]byte(user.Tables), &user)

	var posts []*User

	err := cur.All(ctx, &posts)

	// if jsonErr != nil {
	// 	log.Fatal("jsonErr:::", jsonErr)
	// }

	if err != nil {
		log.Fatal(err)
	}

	return posts[0], nil

}
