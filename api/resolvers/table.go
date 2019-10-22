package resolvers

import (
	"context"
	"time"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDatetimeFromId(id primitive.ObjectID) primitive.DateTime {
	return primitive.NewDateTimeFromTime(id.Timestamp())
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
		Title:        title,
		Slug:         slug,
		Created:      dateTime,
		LastModified: dateTime,
		Subjects:     []*models.Subject{},
		Classes:      []*models.Class{},
		Teachers:     []*models.Teacher{},
	}

	filter := bson.M{"username": auth.Username}
	update := bson.M{"$push": bson.M{"tables": table}}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

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
		"username":   auth.Username,
		"tables._id": id,
	}

	update := bson.M{
		"$set": bson.D{
			{"tables.$.title", title},
			{"tables.$.slug", slug},
			{"tables.$.lastModified", curTime},
		},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &table, nil
}

func Table(ctx context.Context, slug string) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var user *models.User

	filter := bson.M{
		"username":    auth.Username,
		"tables.slug": slug,
	}

	err := DB.Collection("users").FindOne(ctx, filter).Decode(&user)

	helpers.SetWorkloadAmountForTeachers(user.Tables[0].Teachers)

	if err != nil {
		return nil, err
	}

	return user.Tables[0], nil
}

func DeleteTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	table := &models.Table{
		ID: id,
	}

	filter := bson.M{"username": auth.Username}

	update := bson.D{
		{"$pull", bson.D{
			{"tables", bson.D{
				{"_id", id},
			}},
		}},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return table, nil
}

func findTableById(tables []*models.Table, id primitive.ObjectID) *models.Table {
	for i := 0; i < len(tables); i++ {
		if tables[i].ID == id {
			return tables[i]
		}
	}

	return nil
}

func DuplicateTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	usersCollection := DB.Collection("users")

	var user *models.User

	filter := bson.M{"username": auth.Username}

	usersCollection.FindOne(ctx, filter).Decode(&user)

	originalTable := findTableById(user.Tables, id)

	newTitle := originalTable.Title + " (Copy)"
	newID := primitive.NewObjectID()
	dateTime := GetDatetimeFromId(newID)

	newTable := models.Table{
		ID:           newID,
		Title:        newTitle,
		Slug:         originalTable.Slug + "_copy",
		Subjects:     originalTable.Subjects,
		Teachers:     originalTable.Teachers,
		Classes:      originalTable.Classes,
		Created:      dateTime,
		LastModified: dateTime,
	}

	update := bson.M{"$push": bson.M{"tables": newTable}}

	_, err := usersCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &newTable, nil
}
