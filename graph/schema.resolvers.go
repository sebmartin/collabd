package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/sebmartin/collabd/graph/generated"
	"github.com/sebmartin/collabd/models"
)

type TodoError string

func (e TodoError) Error() string {
	return string(e)
}

// StartSession is the resolver for the startSession field.
func (r *mutationResolver) StartSession(ctx context.Context) (*models.Session, error) {
	// return models.NewSession(r.DB)
	return nil, TodoError("we need to bind a session to a game kernel")
}

// JoinSession is the resolver for the joinSession field.
func (r *mutationResolver) JoinSession(ctx context.Context, name string, code string) (*models.Participant, error) {
	// return models.NewParticipant(r.DB, name, code)
	return nil, TodoError("obtain session from the repo (in memory session)")
}

// Sessions is the resolver for the sessions field.
func (r *queryResolver) Sessions(ctx context.Context) ([]*models.Session, error) {
	db := r.DB
	sessions := []*models.Session{}
	// TODO: this only returns 100, consider adopting the Relay protocol for paging
	result := db.Limit(100).Find(&sessions)
	if result.Error != nil {
		log.Print("Failed to retrieve all sessions")
		return nil, result.Error
	}
	return sessions, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
