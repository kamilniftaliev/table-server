package resolvers

import (
	"context"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Classes(ctx context.Context, tableID primitive.ObjectID) ([]*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var user *models.User

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	err := DB.Collection("users").FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user.Tables[0].Classes, nil
}

func CreateClass(ctx context.Context, title string, isDivisible bool, tableID primitive.ObjectID) (*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()

	class := models.Class{
		ID:          id,
		Title:       title,
		IsDivisible: isDivisible,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{"$push": bson.M{"tables.$.classes": class}}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &class, nil
}

func UpdateClass(
	ctx context.Context,
	id primitive.ObjectID,
	title string,
	isDivisible bool,
	tableID primitive.ObjectID,
) (*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	class := models.Class{
		ID:          id,
		Title:       title,
		IsDivisible: isDivisible,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$set": bson.D{
			{"tables.$.classes.$[class].title", title},
			{"tables.$.classes.$[class].isdivisible", isDivisible},
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"class._id": id}},
	}
	updateOptions := &options.UpdateOptions{}
	updateOptions.SetArrayFilters(arrayFilters)

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update, updateOptions)

	if err != nil {
		return nil, err
	}

	return &class, nil
}

func DeleteClass(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Class, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	class := &models.Class{
		ID: id,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.D{
		{"$pull", bson.D{
			{"tables.$.classes", bson.D{
				{"_id", id},
			}},
		}},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return class, nil
}