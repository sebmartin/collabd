package game

import (
	"fmt"

	"github.com/sebmartin/collabd/models"
	"gorm.io/gorm"
)

type SessionService struct {
	DB           *gorm.DB
	LiveSessions []*models.Session
	// GameRegistry []models.Game
}

// func (r *SessionService) RegisterGame(g models.Game) {
// 	log.Printf("Registered game: %s", g.Name())
// 	r.GameRegistry = append(r.GameRegistry, g)
// }

func (r *SessionService) NewSession(stage models.GameStage) (*models.Session, error) {
	session, err := models.NewSession(r.DB, stage)
	r.LiveSessions = append(r.LiveSessions, session)
	// r.GameRegistry = make([]models.Game, 10)
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
	session.ServerEvents[p.ID] = channel
	return session, nil
}
