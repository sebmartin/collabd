package game

import (
	"context"
	"sync"

	"github.com/sebmartin/collabd/models"
)

var (
	gamesMu sync.RWMutex
	games   = make(map[string]Game)
)

// Registers a game making it available from the server. If called twice with the
// same game name, or game is nil, it panics.
func Register(name string, game *Game) {
	gamesMu.Lock()
	defer gamesMu.Unlock()

	if game == nil {
		panic("Game Registry: attempted to register nil game for name " + name)
	}
	if _, other := games[name]; other {
		panic("Game Registry: a game is already registered with the name " + name)
	}
	games[name] = game
}

func NewSession(ctx context.Context, name string) (*models.Session, error) {
	panic("// TODO implement")
}

func SessionForCode(ctx context.Context, code string) (*models.Session, error) {
	panic("// TODO implement")
}

func JoinSession(ctx context.Context, player *models.Player, code string) (*models.Session, error) {
	panic("// TODO implement")
}
