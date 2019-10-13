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

func Teachers(ctx context.Context, tableID primitive.ObjectID) ([]*models.Teacher, error) {
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

	for i := 0; i < len(user.Tables[0].Teachers); i++ {
		workload := user.Tables[0].Teachers[i].Workload
		totalHours := 0

		for j := 0; j < len(workload); j++ {
			totalHours += *workload[j].Hours
		}
		user.Tables[0].Teachers[i].WorkloadAmount = totalHours
	}

	if err != nil {
		return nil, err
	}

	return user.Tables[0].Teachers, nil
}

func CreateTeacher(ctx context.Context, name string, tableID primitive.ObjectID, slug string) (*models.Teacher, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()

	teacher := models.Teacher{
		ID:       id,
		Name:     name,
		Slug:     slug,
		Workload: []*models.Workload{},
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$push": bson.M{"tables.0.teachers": teacher},
		"$set":  bson.M{"tables.0.lastModified": primitive.NewDateTimeFromTime(time.Now())},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

func UpdateTeacher(
	ctx context.Context,
	id primitive.ObjectID,
	name string,
	tableID primitive.ObjectID,
	slug string,
) (*models.Teacher, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	teacher := models.Teacher{
		ID:   id,
		Name: name,
		Slug: slug,
	}

	filter := bson.M{
		"username":            auth.Username,
		"tables._id":          tableID,
		"tables.teachers._id": id,
	}

	update := bson.M{
		"$set": bson.D{
			{"tables.0.teachers.0.name", name},
			{"tables.0.teachers.0.slug", slug},
			{"tables.0.lastModified", primitive.NewDateTimeFromTime(time.Now())},
		},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

func UpdateWorkload(
	ctx context.Context,
	tableID primitive.ObjectID,
	teacherID primitive.ObjectID,
	subjectID primitive.ObjectID,
	classID primitive.ObjectID,
	hours int,
	prevHours int,
) (*models.Workload, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	workload := models.Workload{
		SubjectID: subjectID,
		ClassID:   classID,
		Hours:     &hours,
	}

	oldWorkload := models.Workload{
		SubjectID: subjectID,
		ClassID:   classID,
		Hours:     &prevHours,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	addNewHours := bson.M{
		"$addToSet": bson.M{"tables.$.teachers.$[teacher].workload": workload},
	}
	deleteOldHours := bson.M{
		"$pull": bson.M{"tables.$.teachers.$[teacher].workload": oldWorkload},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"teacher._id": teacherID},
		},
	}
	updateOptions := &options.UpdateOptions{}
	updateOptions.SetArrayFilters(arrayFilters)

	_, err := DB.Collection("users").UpdateOne(ctx, filter, addNewHours, updateOptions)
	DB.Collection("users").UpdateOne(ctx, filter, deleteOldHours, updateOptions)

	if err != nil {
		return nil, err
	}

	return &workload, nil
}

func DeleteTeacher(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Teacher, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	teacher := &models.Teacher{
		ID: id,
	}

	filter := bson.M{
		"username":   auth.Username,
		"tables._id": tableID,
	}

	update := bson.M{
		"$pull": bson.M{
			"tables.0.teachers": bson.M{"_id": id},
		},
		"$set": bson.M{
			"tables.0.lastModified": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := DB.Collection("users").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return teacher, nil
}
