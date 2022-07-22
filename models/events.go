package models

type EventType string

const (
	JoinEventType EventType = "JOIN"
)

type Event interface {
	Type() EventType
}

type JoinEvent struct {
	Participant *Participant
	Channel     chan Event
}

func (e *JoinEvent) Type() EventType {
	return JoinEventType
}
