package models

import (
	"context"
)

type EventType string

const (
	JoinEventType  EventType = "JOIN" // TODO: remove once session_test is refactored to not use this
	ErrorEventType EventType = "ERROR"
)

// TODO choose idiomatic names for these interfaces
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

func NewPlayerEvent(ctx context.Context, eventType EventType, sender *Player) PlayerEvent {
	return &basicPlayerEvent{
		eventType: eventType,
		sender:    sender,
		ctx:       ctx,
	}
}

type basicPlayerEvent struct {
	eventType EventType
	sender    *Player
	ctx       context.Context
}

func (e *basicPlayerEvent) Type() EventType {
	return e.eventType
}

func (e *basicPlayerEvent) Sender() *Player {
	return e.sender
}

func (e *basicPlayerEvent) Context() context.Context {
	return e.ctx
}

func NewServerEvent(eventType EventType) ServerEvent {
	return &basicServerEvent{
		eventType: eventType,
	}
}

type basicServerEvent struct {
	eventType EventType
}

func (e *basicServerEvent) Type() EventType {
	return e.eventType
}

// Event sent from server in response to a player event that generated an error
type ErrorEvent struct {
	ServerEvent
	Error error
}

func NewErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		ServerEvent: NewServerEvent(ErrorEventType),
		Error:       err,
	}
}
