package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Title Title
}
