package resolvers

import (
	"context"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Classes(ctx context.Context, tableID primitive.ObjectID) ([]*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var table *models.Table

	filter := bson.M{"_id": tableID}

	err := DB.Collection("tables").FindOne(ctx, filter).Decode(&table)

	if err != nil {
		return nil, err
	}

	return table.Classes, nil
}

func CreateClass(
	ctx context.Context,
	tableID primitive.ObjectID,
	shift,
	number int,
	sector,
	letter string,
) (*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()

	class := models.Class{
		ID:      id,
		TableID: tableID,
		Shift:   shift,
		Number:  number,
		Sector:  sector,
		Letter:  letter,
	}

	_, err := DB.Collection("classes").InsertOne(ctx, class)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return &class, nil
}

func UpdateClass(
	ctx context.Context,
	id,
	tableID primitive.ObjectID,
	shift,
	number int,
	sector,
	letter string,
) (*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	class := models.Class{
		ID:      id,
		TableID: tableID,
		Shift:   shift,
		Number:  number,
		Sector:  sector,
		Letter:  letter,
	}

	filter := bson.M{
		"_id":     id,
		"tableId": tableID,
	}

	update := bson.M{
		"$set": bson.M{
			"number": number,
			"letter": letter,
			"shift":  shift,
			"sector": sector,
		},
	}

	_, err := DB.Collection("classes").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return &class, nil
}

func DeleteClass(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*primitive.ObjectID, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	filter := bson.M{
		"_id":     id,
		"tableId": tableID,
	}

	DB.Collection("users").FindOneAndDelete(ctx, filter)

	UpdateLastModifiedTime(tableID)

	return &id, nil
}
