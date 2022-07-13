package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sebmartin/collabd/graph/generated"
	"github.com/sebmartin/collabd/models"
)

// Sessions is the resolver for the sessions field.
func (r *queryResolver) Sessions(ctx context.Context) ([]*models.Session, error) {
	db := r.DB
	sessions := []*models.Session{}
	result := db.Limit(100).Find(&sessions)
	fmt.Print(result)
	return sessions, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
