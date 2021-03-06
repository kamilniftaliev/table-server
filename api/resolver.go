package api

import (
	"context"

	// "github.com/graphql-go/graphql"
	"github.com/kamilniftaliev/table-server/api/models"
	"github.com/kamilniftaliev/table-server/api/resolvers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) SignIn(ctx context.Context, username, password string) (*models.Token, error) {
	return resolvers.SignIn(ctx, username, password)
}

// TABLE RESOLVERS
func (r *mutationResolver) CreateTable(ctx context.Context, title, slug string) (*models.Table, error) {
	return resolvers.CreateTable(ctx, title, slug)
}

func (r *mutationResolver) UpdateTable(ctx context.Context, title, slug string, id primitive.ObjectID) (*models.Table, error) {
	return resolvers.UpdateTable(ctx, title, slug, id)
}

func (r *mutationResolver) DeleteTable(ctx context.Context, id primitive.ObjectID) (*primitive.ObjectID, error) {
	return resolvers.DeleteTable(ctx, id)
}

func (r *mutationResolver) DuplicateTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	return resolvers.DuplicateTable(ctx, id)
}

// SUBJECT RESOLVERS
// func (r *mutationResolver) CreateSubject(
// 	ctx context.Context,
// 	title string,
// ) (*models.Subject, error) {
// 	return resolvers.CreateSubject(ctx, title, tableID)
// }

// func (r *mutationResolver) UpdateSubject(
// 	ctx context.Context,
// 	id primitive.ObjectID,
// 	title string,
// 	tableID primitive.ObjectID,
// ) (*models.Subject, error) {
// 	return resolvers.UpdateSubject(ctx, id, title, tableID)
// }

// func (r *mutationResolver) DeleteSubject(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Subject, error) {
// 	return resolvers.DeleteSubject(ctx, id, tableID)
// }

// CLASS RESOLVERS
func (r *mutationResolver) CreateClass(
	ctx context.Context,
	tableID primitive.ObjectID,
	shift,
	number int,
	sector,
	letter string,
) (*models.Class, error) {
	return resolvers.CreateClass(
		ctx,
		tableID,
		shift,
		number,
		sector,
		letter,
	)
}

func (r *mutationResolver) UpdateClass(
	ctx context.Context,
	id,
	tableID primitive.ObjectID,
	shift,
	number int,
	sector,
	letter string,
) (*models.Class, error) {
	return resolvers.UpdateClass(
		ctx,
		id,
		tableID,
		shift,
		number,
		sector,
		letter,
	)
}

func (r *mutationResolver) DeleteClass(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*primitive.ObjectID, error) {
	return resolvers.DeleteClass(ctx, id, tableID)
}

// TEACHER RESOLVERS
func (r *mutationResolver) CreateTeacher(
	ctx context.Context,
	tableID primitive.ObjectID,
	name,
	slug string,
) (*models.Teacher, error) {
	return resolvers.CreateTeacher(ctx, tableID, name, slug)
}

func (r *mutationResolver) UpdateTeacher(
	ctx context.Context,
	id,
	tableID primitive.ObjectID,
	name,
	slug string,
) (*models.Teacher, error) {
	return resolvers.UpdateTeacher(ctx, id, tableID, name, slug)
}

func (r *mutationResolver) UpdateWorkload(
	ctx context.Context,
	tableID,
	teacherID,
	subjectID,
	classID primitive.ObjectID,
	hours int,
) (*models.Workload, error) {
	return resolvers.UpdateWorkload(ctx, tableID, teacherID, subjectID, classID, hours)
}

func (r *mutationResolver) UpdateWorkhour(
	ctx context.Context,
	tableID,
	teacherID primitive.ObjectID,
	day,
	hour int,
	everyHour,
	everyDay,
	value bool,
) (*models.Workhour, error) {
	return resolvers.UpdateWorkhour(ctx, tableID, teacherID, day, hour, everyHour, everyDay, value)
}

func (r *mutationResolver) DeleteTeacher(ctx context.Context, id, tableID primitive.ObjectID) (*primitive.ObjectID, error) {
	return resolvers.DeleteTeacher(ctx, id, tableID)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Table(ctx context.Context, slug string) (*models.Table, error) {
	return resolvers.Table(ctx, slug)
}

func (r *queryResolver) User(ctx context.Context) (*models.User, error) {
	return resolvers.GetUser(ctx)
}

func (r *queryResolver) Subjects(ctx context.Context) ([]*models.Subject, error) {
	return resolvers.Subjects(ctx)
}

func (r *queryResolver) Classes(ctx context.Context, tableID primitive.ObjectID) ([]*models.Class, error) {
	return resolvers.Classes(ctx, tableID)
}

func (r *queryResolver) Teachers(ctx context.Context, tableID primitive.ObjectID) ([]*models.Teacher, error) {
	return resolvers.Teachers(ctx, tableID)
}
