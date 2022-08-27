package join_stage

import (
	"context"
	"testing"
	"time"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newJoinGameStage(min uint, max uint) JoinGame {
	return JoinGame{
		MinPlayers: min,
		MaxPlayers: max,
		StartGame: func(players []*models.Player) models.StageRunner {
			return &nextStage{
				players: players,
			}
		},
	}
}

type nextStage struct {
	players []*models.Player
}

func (s *nextStage) Run(playerEvents <-chan models.PlayerEvent) models.StageRunner {
	return nil
}

func newEventChannel() chan models.PlayerEvent {
	return make(chan models.PlayerEvent, 1000)
}

func newPlayer(name string) *models.Player {
	return &models.Player{
		Name:         name,
		ServerEvents: make(chan models.ServerEvent, 10),
	}
}

func flushServerEvents(events <-chan models.ServerEvent) []models.ServerEvent {
	all := make([]models.ServerEvent, 0, 20)
	for {
		select {
		case event := <-events:
			all = append(all, event)
		case <-time.After(250 * time.Millisecond):
			return all
		}
	}
}

func TestJoinGame_SendsAckEventToPlayers(t *testing.T) {
	stage := newJoinGameStage(1, 3)
	events := newEventChannel()
	go stage.Run(events)

	players := []*models.Player{
		newPlayer("Annie"),
		newPlayer("Steve"),
	}

	for _, p := range players {
		events <- NewJoinEvent(context.Background(), p)
		select {
		case serverEvent := <-p.ServerEvents:
			if errorEvent, ok := serverEvent.(*models.ErrorEvent); ok {
				assert.Failf(t, "Unexpected error event", errorEvent.Error.Error())
				continue
			}
			assert.IsTypef(t, &DidJoinEvent{}, serverEvent, "joingin player: %s", p.Name)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Did not receive a join event acknowledgement from server", "joining player: %s", p.Name)
		}
	}
}

func TestJoinGame_TooManyPlayers(t *testing.T) {
	stage := newJoinGameStage(1, 3)
	events := newEventChannel()
	go stage.Run(events)

	players := []*models.Player{
		newPlayer("Annie"),
		newPlayer("Steve"),
		newPlayer("Joan"),
		newPlayer("Xavier"),
	}

	for _, p := range players {
		events <- NewJoinEvent(context.Background(), p)
	}

	select {
	case event := <-players[3].ServerEvents:
		require.IsType(t, &models.ErrorEvent{}, event, "Expected %s's join request to return an error", players[3].Name)
		errorEvent := event.(*models.ErrorEvent)
		assert.ErrorContains(t, errorEvent.Error, "maximum player count of 3 has already been reached")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Did not receive a response from server")
	}
}

func TestStartGame_NotEnoughPlayers(t *testing.T) {
	stage := newJoinGameStage(2, 3)
	events := newEventChannel()
	go stage.Run(events)

	player := newPlayer("Annie")
	events <- NewJoinEvent(context.Background(), player)
	events <- NewStartEvent(context.Background(), player)

	serverEvents := flushServerEvents(player.ServerEvents)
	require.Len(t, serverEvents, 2)
	assert.Equal(t, DidJoinEventType, serverEvents[0].Type())
	assert.Equal(t, models.ErrorEventType, serverEvents[1].Type())
	assert.ErrorContains(t, serverEvents[1].(*models.ErrorEvent).Error,
		"only 1 player(s) have joined, a minimum of 2 are required before the game can be started",
	)
}

func TestStartGame_ServerEvents(t *testing.T) {
	stage := newJoinGameStage(2, 3)
	events := newEventChannel()
	go stage.Run(events)

	players := []*models.Player{
		newPlayer("Annie"),
		newPlayer("Steve"),
	}

	for _, p := range players {
		events <- NewJoinEvent(context.Background(), p)
	}
	events <- NewStartEvent(context.Background(), players[0])

	for _, p := range players {
		serverEvents := flushServerEvents(p.ServerEvents)
		require.Len(t, serverEvents, 2)
		assert.Equal(t, DidJoinEventType, serverEvents[0].Type(), "player: %s", p.Name)
		assert.Equal(t, DidStartEventType, serverEvents[1].Type(), "player: %s", p.Name)
	}
}

func TestStartGame_NextStage(t *testing.T) {
	stage := newJoinGameStage(4, 4)
	events := newEventChannel()
	var theNextStage models.StageRunner
	go func() {
		theNextStage = stage.Run(events)
	}()

	players := []*models.Player{
		newPlayer("Annie"),
		newPlayer("Steve"),
		newPlayer("Joan"),
		newPlayer("Xavier"),
	}

	for _, p := range players {
		events <- NewJoinEvent(context.Background(), p)
	}
	events <- NewStartEvent(context.Background(), players[0])

	assert.Eventually(t, func() bool {
		return theNextStage != nil
	}, 2*time.Second, 10*time.Millisecond, "Next stage was not returned on time")
	assert.IsType(t, &nextStage{}, theNextStage)
	assert.Len(t, theNextStage.(*nextStage).players, len(players), "The next stage should have received a reference to all players")
}
