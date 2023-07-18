package connect4

import (
	"context"

	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/game/join_stage"
	"github.com/sebmartin/collabd/models"
)

func Register() {
	game.Register("Connect4", func(ctx context.Context) (models.GameDescriber, error) {
		return models.NewGame(
			"Connect 4",
			newInitialStage(),
		), nil
	})
}

func newInitialStage() models.StageRunner {
	return &join_stage.JoinGame{
		MinPlayers: 2,
		MaxPlayers: 2,
		StartGame:  newMainStage,
	}
}

func newMainStage(players []*models.Player) models.StageRunner {
	if len(players) != 2 {
		panic("Connect 4 requires exactly two players")
	}

	return &mainStage{
		players: [2]*models.Player{
			players[0], players[1],
		},
		activePlayer: players[0], // TODO Randomize?
		board:        Board{},
	}
}
