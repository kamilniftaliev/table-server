package resolvers

import (
	"context"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDatetimeFromId(id primitive.ObjectID) string {
	return id.Timestamp().Format("02.01.2006 15:04")
}

func CreateTable(ctx context.Context, title string) (*models.Table, error) {
	userAuth, noAuth := helpers.GetAuthFromContext(ctx)

	if noAuth != nil {
		return nil, noAuth
	}

	id := primitive.NewObjectID()
	dateTime := GetDatetimeFromId(id)

	table := models.Table{
		ID:         id,
		Title:      title,
		Created:    dateTime,
		LastEdited: dateTime,
	}

	filter := bson.M{"username": userAuth.Username}
	update := bson.M{"$push": bson.M{"tables": table}}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &table, nil
}

func DeleteTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	userAuth, noAuth := helpers.GetAuthFromContext(ctx)

	if noAuth != nil {
		return nil, noAuth
	}

	table := &models.Table{
		ID: id,
	}

	filter := bson.M{"username": userAuth.Username}

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
	userAuth, noAuth := helpers.GetAuthFromContext(ctx)

	if noAuth != nil {
		return nil, noAuth
	}

	usersCollection := DB.Collection("users")

	var user *models.User

	filter := bson.M{"username": userAuth.Username}

	usersCollection.FindOne(ctx, filter).Decode(&user)

	originalTable := findTableById(user.Tables, id)

	newTitle := originalTable.Title + " Copy"
	newId := primitive.NewObjectID()
	dateTime := GetDatetimeFromId(newId)

	newTable := models.Table{
		ID:         newId,
		Title:      newTitle,
		Created:    dateTime,
		LastEdited: dateTime,
	}

	update := bson.M{"$push": bson.M{"tables": newTable}}

	_, err := usersCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &newTable, nil
}
