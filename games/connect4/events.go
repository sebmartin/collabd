package connect4

import (
	"context"

	"github.com/sebmartin/collabd/models"
)

const (
	PlayerTurnEventType   = models.EventType("PLAYER_TURN")
	DropPieceEventType    = models.EventType("DROP_PIECE")
	DidDropPieceEventType = models.EventType("DID_DROP_PIECE")
	DidWinEventType       = models.EventType("DID_WIN")
)

type PlayerTurnEvent struct {
	models.ServerEvent

	activePlayer *models.Player
}

func NewPlayerTurnEvent(activePlayer *models.Player) *PlayerTurnEvent {
	return &PlayerTurnEvent{
		ServerEvent:  models.NewServerEvent(PlayerTurnEventType),
		activePlayer: activePlayer,
	}
}

type DropPieceEvent struct {
	models.PlayerEvent

	Slot uint
}

func NewDropPieceEvent(ctx context.Context, sender *models.Player, slot uint) *DropPieceEvent {
	return &DropPieceEvent{
		PlayerEvent: models.NewPlayerEvent(ctx, DropPieceEventType, sender),
		Slot:        slot,
	}
}

type DidDropPieceEvent struct {
	models.ServerEvent

	Piece Piece
	Slot  uint
	Row   uint
}

func NewDidDropPieceEvent(piece Piece, slot uint, row uint) *DidDropPieceEvent {
	return &DidDropPieceEvent{
		ServerEvent: models.NewServerEvent(DidDropPieceEventType),
		Piece:       piece,
		Slot:        slot,
		Row:         row,
	}
}

type DidWinGame struct {
	models.ServerEvent

	Winner models.Player
	Board  Board
}

func NewDidWinGame(winner *models.Player, board *Board) *DidWinGame {
	return &DidWinGame{
		ServerEvent: models.NewServerEvent(DidWinEventType),
		Winner:      *winner,
		Board:       *board,
	}
}
