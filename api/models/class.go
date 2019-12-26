package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Class struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Number int                `json:"number"`
	Letter string             `json:"letter"`
	Shift  int                `json:"shift"`
	Sector string             `json:"sector"`
}
