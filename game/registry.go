package game

import (
	"context"
	"fmt"
	"sync"

	"github.com/sebmartin/collabd/models"
)

var (
	gameRegistryMu sync.RWMutex
	gameRegistry   = make(map[string]func(context.Context) (models.GameInitializer, error))
)

// Registers a game making it available from the server. If called twice with the
// same game name, or game is nil, it panics.
func Register(name string, game func(ctx context.Context) (models.GameInitializer, error)) {
	gameRegistryMu.Lock()
	defer gameRegistryMu.Unlock()

	if game == nil {
		panic("Game Registry: attempted to register nil game for name " + name)
	}
	if _, other := gameRegistry[name]; other {
		panic("Game Registry: a game is already registered with the name " + name)
	}
	gameRegistry[name] = game
}

func NewGame(name string, ctx context.Context) (models.GameInitializer, error) {
	gameRegistryMu.RLock()
	defer gameRegistryMu.RUnlock()

	gameInit, found := gameRegistry[name]
	if !found {
		return nil, fmt.Errorf("failed to create session, unknown game: %s", name)
	}
	return gameInit(ctx)
}
