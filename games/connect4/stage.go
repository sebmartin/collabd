package connect4

import (
	"fmt"

	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/models"
)

type mainStage struct {
	players      [2]*models.Player
	activePlayer *models.Player
	board        Board
}

func (s *mainStage) Run(playerEvents <-chan models.PlayerEvent) models.StageRunner {
	game.Broadcast(s.players[:], NewPlayerTurnEvent(s.activePlayer))

	for {
		event := <-playerEvents
		switch event := event.(type) {
		case DropPieceEvent:
			player := event.Sender()
			if player.ID != s.activePlayer.ID {
				player.ServerEvents <- models.NewErrorEvent(
					fmt.Errorf("player attempted to drop piece when not their turn: %s", player.Name),
				)
				continue
			}

			piece, err := s.playerPiece(player)
			if err != nil {
				player.ServerEvents <- models.NewErrorEvent(err)
				continue
			}

			slot := event.Slot
			row, err := s.board.DropPiece(piece, slot)
			if err != nil {
				player.ServerEvents <- models.NewErrorEvent(err)
				continue
			}

			game.Broadcast(s.players[:], NewDidDropPieceEvent(
				piece, slot, row,
			))

			if s.board.AnalyzeMove(slot, row) == GameWon {
				// We have a winner!
				game.Broadcast(s.players[:], NewDidWinGame(
					player, &s.board,
				))
				return nil
			}

			// Next player's turn
			otherPlayer, err := s.otherPlayer(player)
			if err != nil {
				player.ServerEvents <- models.NewErrorEvent(err)
				continue
			}
			game.Broadcast(s.players[:], NewPlayerTurnEvent(
				otherPlayer,
			))
		}
	}
}

func (s *mainStage) playerPiece(player *models.Player) (Piece, error) {
	if player == s.players[0] {
		return Red, nil
	} else if player == s.players[1] {
		return Black, nil
	}
	return Red, fmt.Errorf("unkonwn player: %s", player.Name)
}

func (s *mainStage) otherPlayer(player *models.Player) (*models.Player, error) {
	if player == s.players[0] {
		return s.players[1], nil
	} else if player == s.players[1] {
		return s.players[0], nil
	}
	return nil, fmt.Errorf("unknown player: %s", player.Name)
}
