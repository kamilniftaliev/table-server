package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Table struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Title         string             `json:"title"`
	Slug          string             `json:"slug"`
	Created       primitive.DateTime `json:"created"`
	LastModified  primitive.DateTime `json:"lastModified" bson:"lastModified"`
	SubjectsCount int                `json:"subjectsCount" bson:"subjectsCount,omitempty"`
	Subjects      []*Subject         `json:"subjects"`
	ClassesCount  int                `json:"classesCount" bson:"classesCount,omitempty"`
	Classes       []*Class           `json:"classes"`
	TeachersCount int                `json:"teachersCount" bson:"teachersCount,omitempty"`
	Teachers      []*Teacher         `json:"teachers"`
}
