package models

import (
	"gorm.io/gorm"
)

type Participant struct {
	*gorm.Model

	Name      string
	SessionID int
	Session   Session
}

func NewParticipant(db *gorm.DB, name string, session *Session) (*Participant, error) {
	p := &Participant{
		Name:    name,
		Session: *session,
	}
	result := db.Create(&p)
	if result.Error != nil {
		return nil, result.Error
	}

	session.ParticipantChannels[p.ID] = make(chan Event)

	return p, nil
}
