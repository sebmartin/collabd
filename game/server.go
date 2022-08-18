package game

import (
	"context"
	"fmt"
	"sync"

	"github.com/sebmartin/collabd/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	db         *gorm.DB
	sessionsMu sync.RWMutex
	sessions   []*models.Session
}

func NewServer(driverName string, dsn string) (*Server, error) {
	dialector, err := func() (gorm.Dialector, error) {
		switch driverName {
		case "postgres":
			return postgres.Open(dsn), nil
		case "mysql":
			return mysql.Open(dsn), nil
		case "sqlite":
			return sqlite.Open(dsn), nil
		default:
			return nil, fmt.Errorf("unsupported driver name: %s", driverName)
		}
	}()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = gormDB.AutoMigrate(
		&models.Player{},
		&models.Session{},
	)
	if err != nil {
		return nil, err
	}

	return &Server{
		db: gormDB,
	}, nil
}

func (s *Server) NewSession(ctx context.Context, name string) (*models.Session, error) {
	game, err := NewGame(name, ctx)
	if err != nil {
		return nil, err
	}

	session, err := models.NewSession(s.db, game)
	if err != nil {
		return nil, err
	}
	s.appendSession(session)
	return session, nil
}

func (s *Server) appendSession(session *models.Session) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	if s.sessions == nil {
		s.sessions = make([]*models.Session, 0, 10)
	}
	s.sessions = append(s.sessions, session)
}

func (s *Server) SessionForCode(code string) (*models.Session, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	for _, s := range s.sessions {
		if s.Code == code {
			return s, nil
		}
	}
	return nil, fmt.Errorf(`could not find session with code "%s"`, code)
}

func (s *Server) JoinSession(ctx context.Context, player *models.Player, code string) (*models.Session, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	session, err := s.SessionForCode(code)
	if err != nil {
		return nil, err
	}
	channel, err := session.AddPlayer(s.db, player)
	if err != nil {
		return nil, err
	}
	session.ServerEvents[player.ID] = channel
	return session, nil
}