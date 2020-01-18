package resolvers

import (
	"context"
	"log"
	"time"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDatetimeFromId(id primitive.ObjectID) primitive.DateTime {
	return primitive.NewDateTimeFromTime(id.Timestamp())
}

func UpdateLastModifiedTime(tableID primitive.ObjectID) {
	filter := bson.M{
		"_id": tableID,
	}

	update := bson.M{
		"lastModified": primitive.NewDateTimeFromTime(time.Now()),
	}

	DB.Collection("tables").UpdateOne(context.Background(), filter, update)
}

func CreateTable(ctx context.Context, title, slug string) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()
	dateTime := GetDatetimeFromId(id)

	table := models.Table{
		ID:           id,
		UserID:       auth.UserID,
		Title:        title,
		Slug:         slug,
		Created:      dateTime,
		LastModified: dateTime,
	}

	_, err := DB.Collection("tables").InsertOne(ctx, table)

	if err != nil {
		return nil, err
	}

	return &table, nil
}

func UpdateTable(ctx context.Context, title, slug string, id primitive.ObjectID) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	dateTime := GetDatetimeFromId(id)
	curTime := primitive.NewDateTimeFromTime(time.Now())

	table := models.Table{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Created:      dateTime,
		LastModified: curTime,
	}

	filter := bson.M{
		"_id":    id,
		"userId": auth.UserID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":        title,
			"slug":         slug,
			"lastModified": curTime,
		},
	}

	DB.Collection("tables").UpdateOne(ctx, filter, update)

	return &table, nil
}

func Table(ctx context.Context, slug string) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var tables []*models.Table

	pipeline := []bson.M{
		bson.M{"$match": bson.M{
			"userId": auth.UserID,
			"slug":   slug,
		}},
		bson.M{"$lookup": bson.M{
			"from":         "teachers",
			"localField":   "_id",
			"foreignField": "tableId",
			"as":           "teachers",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "classes",
			"localField":   "_id",
			"foreignField": "tableId",
			"as":           "classes",
		}},
	}

	cursor, err := DB.Collection("tables").Aggregate(ctx, pipeline)

	if err = cursor.All(ctx, &tables); err != nil {
		log.Fatal(err)
	}

	if len(tables) == 0 {
		return nil, err
	}

	return tables[0], nil
}

func DeleteTable(ctx context.Context, id primitive.ObjectID) (*primitive.ObjectID, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	filter := bson.M{
		"_id":    id,
		"userId": auth.UserID,
	}

	tableIDFilter := bson.M{"tableId": id}

	DB.Collection("tables").DeleteOne(ctx, filter)
	DB.Collection("teachers").DeleteMany(ctx, tableIDFilter)
	DB.Collection("classes").DeleteMany(ctx, tableIDFilter)

	return &id, nil
}

func DuplicateTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	tablesCollection := DB.Collection("tables")
	teachersCollection := DB.Collection("teachers")
	classesCollection := DB.Collection("classes")

	var table *models.Table
	var teachers []*models.Teacher
	var classes []*models.Class

	filter := bson.M{
		"_id":    id,
		"userId": auth.UserID,
	}

	tablesCollection.FindOne(ctx, filter).Decode(&table)

	newTitle := table.Title + " Copy"
	newID := primitive.NewObjectID()
	dateTime := GetDatetimeFromId(newID)

	newTable := models.Table{
		ID:           newID,
		Title:        newTitle,
		UserID:       auth.UserID,
		Slug:         table.Slug + "_copy",
		Created:      dateTime,
		LastModified: dateTime,
	}

	_, err := tablesCollection.InsertOne(ctx, newTable)

	if err != nil {
		return nil, err
	}

	tableIDFilter := bson.M{"tableId": id}

	// TEACHERS DUPLICATION
	teachersResults, teachersErr := teachersCollection.Find(ctx, tableIDFilter)

	if teachersErr != nil {
		return nil, teachersErr
	}

	teachersResults.All(ctx, &teachers)
	for i := 0; i < len(teachers); i++ {
		teachers[i].ID = primitive.NewObjectID()
		teachers[i].TableID = newID
	}

	var teachersDocs []interface{}
	for _, t := range teachers {
		teachersDocs = append(teachersDocs, t)
	}

	_, teachersInsertError := teachersCollection.InsertMany(ctx, teachersDocs)

	if teachersInsertError != nil {
		return nil, teachersInsertError
	}

	// Classes DUPLICATION
	classesResults, classesErr := classesCollection.Find(ctx, tableIDFilter)

	if classesErr != nil {
		return nil, classesErr
	}

	classesResults.All(ctx, &classes)
	for i := 0; i < len(classes); i++ {
		classes[i].ID = primitive.NewObjectID()
		classes[i].TableID = newID
	}

	var classesDocs []interface{}
	for _, t := range classes {
		classesDocs = append(classesDocs, t)
	}

	_, classesInsertError := classesCollection.InsertMany(ctx, classesDocs)

	if classesInsertError != nil {
		return nil, classesInsertError
	}

	return &newTable, nil
}
