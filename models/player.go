package models

import (
	"gorm.io/gorm"
)

type Player struct {
	*gorm.Model

	Name string
	// SessionID int
	// Session      Session // TODO is this relationship backwards? Can players join more than one game?
	ServerEvents chan ServerEvent `gorm:"-:all"`
}

func NewPlayer(db *gorm.DB, name string) (*Player, error) {
	p := &Player{
		Name: name,
		// Session:      *session,
		ServerEvents: make(chan ServerEvent, ChanBufferSize),
	}
	result := db.Create(p)
	if result.Error != nil {
		return nil, result.Error
	}

	return p, nil
}
