package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Participant struct {
	*gorm.Model

	Name      string
	SessionID int
	Session   Session
}

func NewParticipant(db *gorm.DB, name string, sessionCode string) (*Participant, error) {
	session := Session{}
	result := db.Find(&session, "Code = ?", sessionCode)
	if result.RowsAffected == 0 {
		log.Printf("Could not find session with code: %s", sessionCode)
		return nil, fmt.Errorf("session not found: %s", sessionCode)
	}
	return NewParticipantWithSession(db, name, &session)
}

func NewParticipantWithSession(db *gorm.DB, name string, session *Session) (*Participant, error) {
	p := &Participant{
		Name:    name,
		Session: *session,
	}
	result := db.Create(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return p, nil
}
