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

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) (*models.User, error) {
	return resolvers.GetUser(ctx)
}
