package models

import "context"

type EventType string

const (
	JoinEventType EventType = "JOIN"
)

type Event interface {
	Type() EventType
}

type ServerEvent interface {
	Event
}

type PlayerEvent interface {
	Event

	Sender() *Player
	Context() context.Context
}

type BasicPlayerEvent struct {
	sender *Player
	ctx    context.Context
}

func (e *BasicPlayerEvent) Type() EventType {
	return "?"
}

func (e *BasicPlayerEvent) Sender() *Player {
	return e.sender
}

func (e *BasicPlayerEvent) Context() context.Context {
	return e.ctx
}

func NewJoinEvent(ctx context.Context, player *Player, playerChannel chan<- ServerEvent) *JoinEvent {
	return &JoinEvent{
		PlayerEvent: &BasicPlayerEvent{
			sender: player,
			ctx:    ctx,
		},
		Channel: playerChannel,
	}
}

type JoinEvent struct {
	PlayerEvent

	Channel chan<- ServerEvent
}

func (e JoinEvent) Type() EventType {
	return JoinEventType
}