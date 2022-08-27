package join_stage

import (
	"context"

	"github.com/sebmartin/collabd/models"
)

const (
	JoinEventType     models.EventType = "JOIN"
	DidJoinEventType  models.EventType = "DID_JOIN"
	StartEventType    models.EventType = "START"
	DidStartEventType models.EventType = "DID_START"
)

// Send a JoinEvent to add player to the game. A join request can be refused if the game has already started
// or the maximum number of players has been reached. The event's Sender() is assumed to be the joining player.
type JoinEvent struct {
	models.PlayerEvent
}

func NewJoinEvent(ctx context.Context, sender *models.Player) *JoinEvent {
	return &JoinEvent{
		PlayerEvent: models.NewPlayerEvent(ctx, JoinEventType, sender),
	}
}

type DidJoinEvent struct {
	models.ServerEvent

	Player *models.Player
}

func NewDidJoinEvent(player *models.Player) *DidJoinEvent {
	return &DidJoinEvent{
		ServerEvent: models.NewServerEvent(DidJoinEventType),
		Player:      player,
	}
}

// Send a StartEvent to indicate that all players have joined and the game is ready to begin
type StartEvent struct {
	models.PlayerEvent
}

func NewStartEvent(ctx context.Context, sender *models.Player) *StartEvent {
	return &StartEvent{
		PlayerEvent: models.NewPlayerEvent(ctx, StartEventType, sender),
	}
}

type DidStartEvent struct {
	models.ServerEvent
}

func NewDidStartEvent(players []*models.Player) *DidStartEvent {
	return &DidStartEvent{
		ServerEvent: models.NewServerEvent(DidStartEventType),
	}
}
