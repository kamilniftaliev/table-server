package resolvers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SignIn(ctx context.Context, username, password string) (*models.Token, error) {
	auth := helpers.GetAuth(ctx)

	if len(auth.Username) > 1 {
		token := models.Token{
			Token: auth.Token,
		}
		return &token, nil
	}

	var user models.User

	filter := bson.M{"username": username}
	response := DB.Collection("users").FindOne(ctx, filter).Decode(&user)

	// If didn't find any user
	if response != nil || !helpers.ComparePassword(password, user.Password) {
		return nil, errors.New("errorPassword")
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	// expiresAt := time.Now().Add(time.Second * 60).Unix()

	token := &models.Token{
		Token:     helpers.JwtCreate(user.Username, expiresAt),
		ExpiresAt: int(expiresAt),
	}

	return token, nil
}

func GetUser(ctx context.Context) (*models.User, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	usersCollection := DB.Collection("users")

	var user *models.User

	filter := bson.M{"username": auth.Username}

	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}
