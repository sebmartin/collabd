package connect4

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestMainStage(db *gorm.DB) (*mainStage, chan models.PlayerEvent) {
	player1, player2 := newTestPlayer(db, "Alice"), newTestPlayer(db, "Benny")

	events := make(chan models.PlayerEvent, 100)
	stage := newMainStage([]*models.Player{player1, player2})
	go stage.Run(events)

	return stage.(*mainStage), events
}

func newTestPlayer(db *gorm.DB, name string) *models.Player {
	player, _ := models.NewPlayer(db, name)
	return player
}

func playPiece(stage *mainStage, events chan models.PlayerEvent, player *models.Player, slot uint) {
	events <- NewDropPieceEvent(context.Background(), player, slot)
}

func flushServerEvents(t *testing.T, events <-chan models.ServerEvent, count int) []models.ServerEvent {
	all := make([]models.ServerEvent, 0, 20)
	for {
		select {
		case event := <-events:
			all = append(all, event)
			if len(all) >= count {
				return all
			}
		case <-time.After(250 * time.Millisecond):
			assert.Fail(t, "Timeout waiting for server event", "Events captured so far: %s", all)
			return nil
		}
	}
}

func assertServerEvents(t *testing.T, player *models.Player, expectedEvents []models.ServerEvent) {
	actualEvents := flushServerEvents(t, player.ServerEvents, len(expectedEvents))
	if actualEvents == nil {
		return
	}

	require.Equal(t, len(expectedEvents), len(actualEvents), "Unexpected server event count, events = %s", actualEvents)
	for i, _ := range expectedEvents {
		assert.Equal(t, expectedEvents[i], actualEvents[i], "Unexpected event at index: %d", i)
	}
}

func Test_mainStage_WinGame(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, events := newTestMainStage(db)
	player1 := stage.players[0]
	player2 := stage.players[1]

	// Players are both notified of who's turn it is next as the game starts
	serverEvents := []models.ServerEvent{
		NewPlayerTurnEvent(player1),
	}
	assertServerEvents(t, player1, serverEvents)
	assertServerEvents(t, player2, serverEvents)

	// player1 plays:
	// | . . . . . . . |
	// | . . . . . . . | (row 4)
	// | . R . . . . . | (row 5)
	playPiece(stage, events, player1, 1)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Red, 1, 5),
		NewPlayerTurnEvent(player2),
	}
	assertServerEvents(t, player1, serverEvents)
	assertServerEvents(t, player2, serverEvents)

	// player2 plays:
	// | . . . . . . . |
	// | . B . . . . . | (row 4)
	// | . R . . . . . | (row 5)
	playPiece(stage, events, player2, 1)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Black, 1, 4),
		NewPlayerTurnEvent(player1),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)

	// player1 plays:
	// | . . . . . . . |
	// | . B . . . . . | (row 4)
	// | . R R . . . . | (row 5)
	playPiece(stage, events, player1, 2)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Red, 2, 5),
		NewPlayerTurnEvent(player2),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)

	// player2 plays:
	// | . . . . . . . |
	// | . B B . . . . | (row 4)
	// | . R R . . . . | (row 5)
	playPiece(stage, events, player2, 2)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Black, 2, 4),
		NewPlayerTurnEvent(player1),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)

	// player1 plays:
	// | . . . . . . . |
	// | . B B . . . . | (row 4)
	// | . R R R . . . | (row 5)
	playPiece(stage, events, player1, 3)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Red, 3, 5),
		NewPlayerTurnEvent(player2),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)

	// player2 plays:
	// | . . . . . . . |
	// | . B B B . . . | (row 4)
	// | . R R R . . . | (row 5)
	playPiece(stage, events, player2, 3)
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Black, 3, 4),
		NewPlayerTurnEvent(player1),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)

	// player1 WINS!
	// | . . . . . . . |
	// | . B B B . . . | (row 4)
	// | . R R R R . . | (row 5)
	playPiece(stage, events, player1, 4)
	R, B := Red, Black
	serverEvents = []models.ServerEvent{
		NewDidDropPieceEvent(Red, 4, 5),
		NewDidWinGame(player1, &Board{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, B, B, B, 0, 0, 0},
			{0, R, R, R, R, 0, 0},
		}),
	}
	assertServerEvents(t, player2, serverEvents)
	assertServerEvents(t, player1, serverEvents)
}

func Test_mainStage_PlayedOutOfTurn(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, events := newTestMainStage(db)
	player1 := stage.players[0]
	player2 := stage.players[1]

	playPiece(stage, events, player1, 1)
	playPiece(stage, events, player1, 2)
	assertServerEvents(t, player1, []models.ServerEvent{
		NewPlayerTurnEvent(player1),
		NewDidDropPieceEvent(Red, 1, 5),
		NewPlayerTurnEvent(player2),
		models.NewErrorEvent(fmt.Errorf("player attempted to drop piece when not their turn: Alice")),
	})
	assertServerEvents(t, player2, []models.ServerEvent{
		NewPlayerTurnEvent(player1),
		NewDidDropPieceEvent(Red, 1, 5),
		NewPlayerTurnEvent(player2),
	})
}

func Test_mainStage_ImposterPlayerEvent(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, events := newTestMainStage(db)
	imposter := newTestPlayer(db, "Imposter")

	playPiece(stage, events, imposter, 1)
	assertServerEvents(t, imposter, []models.ServerEvent{
		models.NewErrorEvent(fmt.Errorf("unknown player: Imposter")),
	})
}

func Test_mainStage_InvalidSlot(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, events := newTestMainStage(db)
	player1 := stage.players[0]

	playPiece(stage, events, player1, MaxColumns)
	assertServerEvents(t, player1, []models.ServerEvent{
		NewPlayerTurnEvent(player1),
		models.NewErrorEvent(fmt.Errorf("slot 7 exceeds the slot maximum of 6")),
	})
}

func Test_mainStage_otherPlayer(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, _ := newTestMainStage(db)
	player1 := stage.players[0]
	player2 := stage.players[1]

	result, err := stage.otherPlayer(player1)
	assert.Equal(t, result, player2)
	assert.Nil(t, err)

	result, err = stage.otherPlayer(player2)
	assert.Equal(t, result, player1)
	assert.Nil(t, err)
}

func Test_mainStage_otherPlayer_unknownPlayer(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, _ := newTestMainStage(db)
	imposter := newTestPlayer(db, "Imposter")

	result, err := stage.otherPlayer(imposter)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "unknown player: Imposter")
}

func Test_mainStage_playerPiece(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, _ := newTestMainStage(db)
	player1 := stage.players[0]
	player2 := stage.players[1]

	piece, err := stage.playerPiece(player1)
	assert.Equal(t, piece, Red)
	assert.Nil(t, err)

	piece, err = stage.playerPiece(player2)
	assert.Equal(t, piece, Black)
	assert.Nil(t, err)
}

func Test_mainStage_playerPiece_unknownPlayer(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	stage, _ := newTestMainStage(db)
	imposter := newTestPlayer(db, "Imposter")

	_, err := stage.playerPiece(imposter)
	assert.ErrorContains(t, err, "unknown player: Imposter")
}
