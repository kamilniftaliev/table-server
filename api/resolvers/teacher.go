package resolvers

import (
	"context"
	"strconv"

	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/kamilniftaliev/table-server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Teachers(ctx context.Context, tableID primitive.ObjectID) ([]*models.Teacher, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var teachers []*models.Teacher

	filter := bson.M{"tableId": tableID}

	results, err := DB.Collection("teachers").Find(ctx, filter)
	results.All(ctx, &teachers)

	if err != nil {
		return nil, err
	}

	return teachers, nil
}

func CreateTeacher(ctx context.Context, tableID primitive.ObjectID, name, slug string) (*models.Teacher, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	id := primitive.NewObjectID()

	// var workhours1 [][]bool
	// for i := 0; i < 5; i++ {
	// 	arr := append(workhours1, [5]bool{})
	// 	for j := 0; j < 5; j++ {
	// 		arr[i][j] = true
	// 	}
	// }
	workhours := [][]bool{
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
	}

	teacher := models.Teacher{
		ID:        id,
		TableID:   tableID,
		Name:      name,
		Slug:      slug,
		Workload:  []*models.Workload{},
		Workhours: workhours,
	}

	_, err := DB.Collection("teachers").InsertOne(ctx, teacher)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return &teacher, nil
}

func UpdateTeacher(
	ctx context.Context,
	id,
	tableID primitive.ObjectID,
	name,
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
		"tableId": tableID,
		"_id":     id,
	}

	update := bson.M{
		"$set": bson.M{
			"name": name,
			"slug": slug,
		},
	}

	_, err := DB.Collection("teachers").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return &teacher, nil
}

func UpdateWorkload(
	ctx context.Context,
	tableID,
	teacherID,
	subjectID,
	classID primitive.ObjectID,
	hours int,
) (*models.Workload, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	var teacher *models.Teacher

	newWorkload := &models.Workload{
		SubjectID: subjectID,
		ClassID:   classID,
		Hours:     hours,
	}

	filter := bson.M{
		"tableId": tableID,
		"_id":     teacherID,
	}

	DB.Collection("teachers").FindOne(ctx, filter).Decode(&teacher)

	workloadIndex := -1

	for i := 0; i < len(teacher.Workload); i++ {
		workload := teacher.Workload[i]
		if workload.SubjectID == subjectID && workload.ClassID == classID {
			workloadIndex = i
		}
	}

	update := bson.M{
		"$set": bson.M{
			"workload." + strconv.Itoa(workloadIndex) + ".hours": hours,
		},
	}

	if workloadIndex == -1 {
		update = bson.M{
			"$push": bson.M{
				"workload": newWorkload,
			},
		}
	}

	_, err := DB.Collection("teachers").UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return newWorkload, nil
}

func UpdateWorkhour(
	ctx context.Context,
	tableID,
	teacherID primitive.ObjectID,
	day,
	hour int,
	everyHour,
	everyDay,
	value bool,
) (*models.Workhour, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	stringDay := strconv.Itoa(day)
	stringHour := strconv.Itoa(hour)

	workhour := models.Workhour{
		Day:       day,
		Hour:      hour,
		Value:     value,
		EveryHour: everyHour,
		EveryDay:  everyDay,
	}

	filter := bson.M{
		"_id":     teacherID,
		"tableId": tableID,
	}

	update := bson.M{
		"$set": bson.M{"workhours." + stringDay + "." + stringHour: value},
	}

	// Switch every day of given hour
	if everyDay {
		update["$set"] = bson.M{"workhours.$[]." + stringHour: value}
	}

	// Switch every hour of given day
	if everyHour {
		if hour == 0 {
			update["$set"] = bson.M{
				"workhours." + stringDay + ".0": value,
				"workhours." + stringDay + ".1": value,
				"workhours." + stringDay + ".2": value,
				"workhours." + stringDay + ".3": value,
				"workhours." + stringDay + ".4": value,
				"workhours." + stringDay + ".5": value,
				"workhours." + stringDay + ".6": value,
				"workhours." + stringDay + ".7": value,
			}
		} else {
			update["$set"] = bson.M{
				"workhours." + stringDay + ".8":  value,
				"workhours." + stringDay + ".9":  value,
				"workhours." + stringDay + ".10": value,
				"workhours." + stringDay + ".11": value,
				"workhours." + stringDay + ".12": value,
				"workhours." + stringDay + ".13": value,
				"workhours." + stringDay + ".14": value,
				"workhours." + stringDay + ".15": value,
			}
		}
	}

	_, err := DB.Collection("teachers").UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	UpdateLastModifiedTime(tableID)

	return &workhour, nil
}

func DeleteTeacher(ctx context.Context, id, tableID primitive.ObjectID) (*primitive.ObjectID, error) {
	auth := helpers.GetAuth(ctx)

	if auth.Error != nil {
		return nil, auth.Error
	}

	filter := bson.M{
		"_id":     id,
		"tableId": tableID,
	}

	DB.Collection("teachers").DeleteOne(ctx, filter)

	UpdateLastModifiedTime(tableID)

	return &id, nil
}
