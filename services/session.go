package services

import (
	"fmt"

	"github.com/sebmartin/collabd/models"
	"gorm.io/gorm"
)

type SessionService struct {
	DB           *gorm.DB
	LiveSessions []*models.Session
}

func (r *SessionService) NewSession(kernel models.GameKernel) (*models.Session, error) {
	session, err := models.NewSession(r.DB, kernel)
	r.LiveSessions = append(r.LiveSessions, session)
	return session, err
}

func (r *SessionService) SessionForCode(code string) (*models.Session, error) {
	for _, s := range r.LiveSessions {
		if s.Code == code {
			return s, nil
		}
	}
	return nil, fmt.Errorf(`could not find session with code "%s"`, code)
}

func (r *SessionService) JoinSession(p *models.Player, code string) (*models.Session, error) {
	session, err := r.SessionForCode(code)
	if err != nil {
		return nil, err
	}
	channel, err := session.AddPlayer(r.DB, p)
	if err != nil {
		return nil, err
	}
	session.PlayerChannels[p.ID] = channel
	return session, nil
}
