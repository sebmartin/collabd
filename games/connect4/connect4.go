package connect4

import (
	"context"

	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/game/join_stage"
	"github.com/sebmartin/collabd/models"
)

func init() {
	game.Register("Connect4", func(ctx context.Context) (models.GameDescriber, error) {
		return models.NewGame(
			"Connect 4",
			&join_stage.JoinGame{
				MinPlayers: 2,
				MaxPlayers: 2,
				StartGame:  initialStage,
			},
		), nil
	})
}

func initialStage(players []*models.Player) models.StageRunner {
	if len(players) != 2 {
		panic("Connect 4 requires exactly two players")
	}

	return &stage{
		players: [2]*models.Player{
			players[0], players[1],
		},
		board: Board{},
	}
}
