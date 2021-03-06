package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Class struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	TableID primitive.ObjectID `json:"tableId" bson:"tableId"`
	Number  int
	Letter  string
	Shift   int
	Sector  string

	Teachers int `json:"teachers" bson:"teachers, omitempty"`
	Subjects int `json:"subjects" bson:"subjects, omitempty"`
	Lessons  int `json:"lessons" bson:"lessons, omitempty"`
}
