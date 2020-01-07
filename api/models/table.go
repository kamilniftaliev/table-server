package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Table struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	UserID        primitive.ObjectID `json:"userId" bson:"userId"`
	Title         string
	Slug          string
	Created       primitive.DateTime
	LastModified  primitive.DateTime `json:"lastModified" bson:"lastModified"`
	ClassesCount  int                `json:"classesCount" bson:"classesCount,omitempty"`
	Classes       []*Class
	TeachersCount int `json:"teachersCount" bson:"teachersCount,omitempty"`
	Teachers      []*Teacher
}
