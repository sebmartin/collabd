package graph

import (
	"github.com/sebmartin/collabd/game"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         *gorm.DB
	GameServer *game.Server
}
