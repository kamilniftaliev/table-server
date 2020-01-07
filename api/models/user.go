package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string
	Username string
	Password string
	Tables   []*Table
}

type Auth struct {
	UserID    primitive.ObjectID `json:"id" bson:"_id"`
	Roles     []string
	IPAddress string
	Token     string
	Error     error
}
