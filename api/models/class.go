package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Class struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title"`
	IsDivisible bool               `json:"isDivisible"`
}
