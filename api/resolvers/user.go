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

	var user *models.User
	var tables []*models.Table

	filter := bson.M{"_id": auth.UserID}

	err := DB.Collection("users").FindOne(ctx, filter).Decode(&user)

	userTableFilter := bson.M{
		"userId": auth.UserID,
	}

	results, err := DB.Collection("tables").Find(ctx, userTableFilter)

	results.All(ctx, &tables)

	for i := 0; i < len(tables); i++ {
		tableStatsFilter := bson.M{"tableId": tables[i].ID}

		teachersCount, _ := DB.Collection("teachers").CountDocuments(ctx, tableStatsFilter)
		classesCount, _ := DB.Collection("classes").CountDocuments(ctx, tableStatsFilter)

		tables[i].TeachersCount = teachersCount
		tables[i].ClassesCount = classesCount
	}

	user.Tables = tables

	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}
