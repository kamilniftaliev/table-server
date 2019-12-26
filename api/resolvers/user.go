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

	if !auth.UserID.IsZero() {
		token := models.Token{
			Token: auth.Token,
		}
		return &token, nil
	}

	var user models.User

	filter := bson.M{"username": username}
	err := DB.Collection("users").FindOne(ctx, filter).Decode(&user)

	// If didn't find any user
	if err != nil || !helpers.ComparePassword(password, user.Password) {
		return nil, errors.New("errorPassword")
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()

	token := &models.Token{
		Token:     helpers.JwtCreate(user.ID, expiresAt),
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

	filter := bson.M{"_id": auth.UserID}

	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	// log.Println(user.Tables)

	if len(user.Tables) > 0 {
		for i := 0; i < len(user.Tables); i++ {
			// user.Tables[i].SubjectsCount = len(user.Tables[i].Subjects)
			user.Tables[i].TeachersCount = len(user.Tables[i].Teachers)
			user.Tables[i].ClassesCount = len(user.Tables[i].Classes)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}
