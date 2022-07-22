package repositories

import (
	"fmt"

	"github.com/sebmartin/collabd/models"
	"gorm.io/gorm"
)

type SessionRepo struct {
	DB           *gorm.DB
	LiveSessions []*models.Session
}

// TODO: rename to SessionRepo everywhere
func (r *SessionRepo) NewSession(kernel models.GameKernel) (*models.Session, error) {
	session, err := models.NewSession(r.DB, kernel)
	r.LiveSessions = append(r.LiveSessions, session)
	return session, err
}

func (r *SessionRepo) SessionForCode(code string) (*models.Session, error) {
	for _, s := range r.LiveSessions {
		if s.Code == code {
			return s, nil
		}
	}
	return nil, fmt.Errorf(`could not find session with code "%s"`, code)
}

func (r *SessionRepo) AddParticipantToSession(p *models.Participant, code string) (*models.Session, error) {
	session, err := r.SessionForCode(code)
	if err != nil {
		return nil, err
	}
	channel, err := session.AddParticipant(r.DB, p)
	if err != nil {
		return nil, err
	}
	session.ParticipantChannels[p.ID] = channel
	return session, nil
}
