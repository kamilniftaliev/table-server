package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workhour struct {
	Day   *string `json:"day"`
	Hour  *string `json:"hour"`
	Value *bool   `json:"value"`
}

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
	// Workhours       [10][10]bool       `json:"workhours"`
	Workhours       [][]bool `json:"workhours"`
	WorkhoursAmount int      `json:"workhoursAmount" bson:"workhoursAmount,omitempty"`
}
