package resolvers

import (
	"context"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Subjects(ctx context.Context) ([]*models.Subject, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var subjects []*models.Subject

	filter := bson.M{}

	results, err := DB.Collection("subjects").Find(ctx, filter)

	if err == nil {
		results.All(ctx, &subjects)
	} else {
		return nil, err
	}

	return subjects, nil
}

// func CreateSubject(ctx context.Context, title string, tableID primitive.ObjectID) (*models.Subject, error) {
// 	auth := helpers.GetAuth(ctx)

// 	if auth.Error != nil {
// 		return nil, auth.Error
// 	}

// 	id := primitive.NewObjectID()

// 	subject := models.Subject{
// 		ID:    id,
// 		Title: title,
// 	}

// 	filter := bson.M{
// 		"username":   auth.UserID,
// 		"tables._id": tableID,
// 	}

// 	update := bson.M{
// 		"$push": bson.M{"tables.$.subjects": subject},
// 		"$set":  bson.M{"tables.$.lastModified": primitive.NewDateTimeFromTime(time.Now())},
// 	}

// 	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &subject, nil
// }

// func UpdateSubject(
// 	ctx context.Context,
// 	id primitive.ObjectID,
// 	title string,
// 	tableID primitive.ObjectID,
// ) (*models.Subject, error) {
// 	auth := helpers.GetAuth(ctx)

// 	if auth.Error != nil {
// 		return nil, auth.Error
// 	}

// 	subject := models.Subject{
// 		ID:    id,
// 		Title: title,
// 	}

// 	filter := bson.M{
// 		"username":   auth.UserID,
// 		"tables._id": tableID,
// 	}

// 	update := bson.M{
// 		"$set": bson.D{
// 			{"tables.$.subjects.$[subject].title", title},
// 			{"tables.$.lastModified", primitive.NewDateTimeFromTime(time.Now())},
// 		},
// 	}

// 	arrayFilters := options.ArrayFilters{
// 		Filters: []interface{}{bson.M{"subject._id": id}},
// 	}
// 	updateOptions := &options.UpdateOptions{}
// 	updateOptions.SetArrayFilters(arrayFilters)

// 	_, err := DB.Collection("users").UpdateOne(ctx, filter, update, updateOptions)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &subject, nil
// }

// func DeleteSubject(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Subject, error) {
// 	auth := helpers.GetAuth(ctx)

// 	if auth.Error != nil {
// 		return nil, auth.Error
// 	}

// 	subject := &models.Subject{
// 		ID: id,
// 	}

// 	filter := bson.M{
// 		"username":   auth.UserID,
// 		"tables._id": tableID,
// 	}

// 	update := bson.M{
// 		"$pull": bson.M{
// 			"tables.$.subjects": bson.M{"_id": id},
// 		},
// 		"$set": bson.M{
// 			"tables.$.lastModified": primitive.NewDateTimeFromTime(time.Now()),
// 		},
// 	}

// 	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return subject, nil
// }
