package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sebmartin/collabd/graph/generated"
	"github.com/sebmartin/collabd/models"
)

// StartSession is the resolver for the startSession field.
func (r *mutationResolver) StartSession(ctx context.Context, gameName *string) (*models.Session, error) {
	return r.GameServer.NewSession(ctx, gameName)
}

// JoinSession is the resolver for the joinSession field.
func (r *mutationResolver) JoinSession(ctx context.Context, name string, code string) (*models.Player, error) {
	// return models.NewPlayer(r.DB, name, code)
	return nil, fmt.Errorf("obtain session from the repo (in memory session)")
}

// Session is the resolver for the session field.
func (r *playerResolver) Session(ctx context.Context, obj *models.Player) (*models.Session, error) {
	panic(fmt.Errorf("not implemented"))
}

// GamesList is the resolver for the gamesList field.
func (r *queryResolver) GamesList(ctx context.Context) ([]string, error) {
	return r.GameServer.GamesList()
}

// Sessions is the resolver for the sessions field.
func (r *queryResolver) Sessions(ctx context.Context) ([]*models.Session, error) {
	sessions := r.GameServer.ActiveSessions()
	if sessions != nil {
		return sessions, nil
	} else {
		return []*models.Session{}, nil
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Player returns generated.PlayerResolver implementation.
func (r *Resolver) Player() generated.PlayerResolver { return &playerResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type playerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
