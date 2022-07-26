package models

import (
	"gorm.io/gorm"
)

type Player struct {
	*gorm.Model

	Name      string
	SessionID int
	Session   Session
}

func NewPlayer(db *gorm.DB, name string, session *Session) (*Player, error) {
	p := &Player{
		Name:    name,
		Session: *session,
	}
	result := db.Create(&p)
	if result.Error != nil {
		return nil, result.Error
	}

	session.PlayerChannels[p.ID] = make(chan Event)

	return p, nil
}
