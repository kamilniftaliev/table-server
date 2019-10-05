package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Table struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Title        string             `json:"title"`
	Slug         string             `json:"slug"`
	Created      primitive.DateTime `json:"created"`
	LastModified primitive.DateTime `json:"lastModified"`
}
