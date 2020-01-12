package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Teacher struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TableID   primitive.ObjectID `json:"tableId" bson:"tableId"`
	Name      string
	Slug      string
	Workload  []*Workload
	Workhours [][]bool

	Subjects        int `json:"subjects" bson:"subjects, omitempty"`
	Classes         int `json:"classes" bson:"classes, omitempty"`
	WorkhoursAmount int `json:"workhoursAmount" bson:"workhoursAmount, omitempty"`
	WorkloadAmount  int `json:"workloadAmount" bson:"workloadAmount, omitempty"`
}
