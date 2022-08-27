package connect4

import "github.com/sebmartin/collabd/models"

type stage struct {
	players [2]*models.Player
	board   Board
}

func (s *stage) Run(<-chan models.PlayerEvent) models.StageRunner {
	return nil // TODO finish game
}
