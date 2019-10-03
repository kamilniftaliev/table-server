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

func FindUsername(ctx context.Context, username, password string) (*models.Token, error) {
	userAuth, _ := helpers.GetAuthFromContext(ctx)

	if len(userAuth.Username) > 1 {
		token := models.Token{
			Token: userAuth.Token,
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

	token := &models.Token{
		Token:     helpers.JwtCreate(user.Username, expiresAt),
		ExpiresAt: int(expiresAt),
	}

	return token, nil
}

func GetUser(ctx context.Context) (*models.User, error) {
	auth, authErr := helpers.GetAuthFromContext(ctx)

	if authErr != nil {
		return nil, authErr
	}

	usersCollection := DB.Collection("users")

	var user *models.User

	filter := bson.M{"username": auth.Username}

	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	// user.ID.Timestamp()
	// for i := 0; i < len(user.Tables); i++ {
	// 	created := GetDatetimeFromId(user.Tables[i].ID)

	// 	user.Tables[i].Created = created
	// 	user.Tables[i].LastEdited = created
	// }

	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}
