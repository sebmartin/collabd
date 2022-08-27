package connect4

import (
	"context"

	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/models"
)

func init() {
	game.Register("Connect4", func(ctx context.Context) (models.GameDescriber, error) {
		return NewGame(), nil
	})
}

const (
	MaxColumns uint = 7
	MaxRows    uint = 6
)

type Connect4 struct {
	models.Game

	Board Board
}

func NewGame() Connect4 {
	return Connect4{
		Game: *models.NewGame(
			"Connect 4",
			nil, // TODO: set initialStage
		),
		Board: [6][7]Piece{},
	}
}
