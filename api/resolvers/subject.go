package resolvers

import (
	"context"
	"time"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Subjects(ctx context.Context, tableID primitive.ObjectID) ([]*models.Subject, error) {
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

	tableIndex := 0

	for i := 0; i < len(user.Tables); i++ {
		if user.Tables[i].ID == tableID {
			tableIndex = i
		}
	}

	if err != nil {
		return nil, err
	}

	return user.Tables[tableIndex].Subjects, nil
}

func CreateSubject(ctx context.Context, title string, tableID primitive.ObjectID) (*models.Subject, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()

	subject := models.Subject{
		ID:    id,
		Title: title,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$push": bson.M{"tables.$.subjects": subject},
		"$set":  bson.M{"tables.$.lastModified": primitive.NewDateTimeFromTime(time.Now())},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &subject, nil
}

func UpdateSubject(
	ctx context.Context,
	id primitive.ObjectID,
	title string,
	tableID primitive.ObjectID,
) (*models.Subject, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	subject := models.Subject{
		ID:    id,
		Title: title,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$set": bson.D{
			{"tables.$.subjects.$[subject].title", title},
			{"tables.$.lastModified", primitive.NewDateTimeFromTime(time.Now())},
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"subject._id": id}},
	}
	updateOptions := &options.UpdateOptions{}
	updateOptions.SetArrayFilters(arrayFilters)

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update, updateOptions)

	if err != nil {
		return nil, err
	}

	return &subject, nil
}

func DeleteSubject(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Subject, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	subject := &models.Subject{
		ID: id,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$pull": bson.M{
			"tables.$.subjects": bson.M{"_id": id},
		},
		"$set": bson.M{
			"tables.$.lastModified": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return subject, nil
}
