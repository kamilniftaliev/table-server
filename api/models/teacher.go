package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workload struct {
	SubjectID primitive.ObjectID `json:"subjectId" bson:"subjectId"`
	ClassID   primitive.ObjectID `json:"classId" bson:"classId"`
	Hours     *int               `json:"hours"`
}

type Teacher struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name"`
	Slug           string             `json:"slug"`
	Workload       []*Workload        `json:"workload"`
	WorkloadAmount int                `json:"workloadAmount" bson:"workloadAmount,omitempty"`
}
