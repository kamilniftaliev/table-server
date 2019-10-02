package api

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func FindUsername(ctx context.Context, username string) (*User, error) {
	cur := DB.Collection("users").FindOne(
		ctx,
		bson.M{"username": username},
	)

	err := cur.Err()

	if err != nil {
		log.Println("Error findUsername: ", err)
		return nil, err
	}

	var user User

	cur.Decode(&user)

	return &user, nil
}
