package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Tables   []*Table           `json:"tables"`
}