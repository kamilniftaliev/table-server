package table

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

	expiresAt := time.Now().Add(time.Hour * 1).Unix()

	token := &Token{
		Token:     JwtCreate(user.Username, expiresAt),
		ExpiresAt: int(expiresAt),
	}

	return token, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
	var users []*User
	usersCollection := DB.Collection("users")

	auth := GetAuthFromContext(ctx)
	log.Println("auth: ", auth)
	cur, err := usersCollection.Find(
		ctx,
		bson.D{},
	)

	// Close the cursor once finished
	defer cur.Close(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var user *User

		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return users, nil

}
