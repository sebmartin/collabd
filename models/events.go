package models

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
}

type BasicPlayerEvent struct {
	sender *Player
}

func (e *BasicPlayerEvent) Type() EventType {
	return "?"
}

func (e *BasicPlayerEvent) Sender() *Player {
	return e.sender
}

func NewJoinEvent(player *Player, playerChannel chan<- ServerEvent) *JoinEvent {
	return &JoinEvent{
		PlayerEvent: &BasicPlayerEvent{
			sender: player,
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
