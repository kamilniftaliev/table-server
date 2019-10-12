package api

import (
	"context"

	"github.com/kamilniftaliev/table-server/api/models"
	"github.com/kamilniftaliev/table-server/api/resolvers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *queryResolver) Table(ctx context.Context, slug string) (*models.Table, error) {
	return resolvers.Table(ctx, slug)
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

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

func (r *mutationResolver) DeleteTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	return resolvers.DeleteTable(ctx, id)
}

func (r *mutationResolver) DuplicateTable(ctx context.Context, id primitive.ObjectID) (*models.Table, error) {
	return resolvers.DuplicateTable(ctx, id)
}

// SUBJECT RESOLVERS
func (r *mutationResolver) CreateSubject(
	ctx context.Context,
	title string,
	isDivisible bool,
	tableID primitive.ObjectID,
) (*models.Subject, error) {
	return resolvers.CreateSubject(ctx, title, isDivisible, tableID)
}

func (r *mutationResolver) UpdateSubject(
	ctx context.Context,
	id primitive.ObjectID,
	title string,
	isDivisible bool,
	tableID primitive.ObjectID,
) (*models.Subject, error) {
	return resolvers.UpdateSubject(ctx, id, title, isDivisible, tableID)
}

func (r *mutationResolver) DeleteSubject(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Subject, error) {
	return resolvers.DeleteSubject(ctx, id, tableID)
}

// CLASS RESOLVERS
func (r *mutationResolver) CreateClass(
	ctx context.Context,
	title string,
	isDivisible bool,
	tableID primitive.ObjectID,
) (*models.Class, error) {
	return resolvers.CreateClass(ctx, title, isDivisible, tableID)
}

func (r *mutationResolver) UpdateClass(
	ctx context.Context,
	id primitive.ObjectID,
	title string,
	isDivisible bool,
	tableID primitive.ObjectID,
) (*models.Class, error) {
	return resolvers.UpdateClass(ctx, id, title, isDivisible, tableID)
}

func (r *mutationResolver) DeleteClass(ctx context.Context, id primitive.ObjectID, tableID primitive.ObjectID) (*models.Class, error) {
	return resolvers.DeleteClass(ctx, id, tableID)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) (*models.User, error) {
	return resolvers.GetUser(ctx)
}

func (r *queryResolver) Subjects(ctx context.Context, tableID primitive.ObjectID) ([]*models.Subject, error) {
	return resolvers.Subjects(ctx, tableID)
}

func (r *queryResolver) Classes(ctx context.Context, tableID primitive.ObjectID) ([]*models.Class, error) {
	return resolvers.Classes(ctx, tableID)
}
