package models

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const (
	SessionCodeLength = 4
	SessionCodeMax    = 456976 // 26^4
)

type Session struct {
	gorm.Model

	Code    string
	Players []Player

	CurrentStage GameStage                 `gorm:"-:all"`
	ServerEvents map[uint]chan ServerEvent `gorm:"-:all"`
	PlayerEvents chan PlayerEvent          `gorm:"-:all"`
}

func (s *Session) AfterCreate(tx *gorm.DB) error {
	initSession(s)
	return nil
}

func (s *Session) AfterFind(tx *gorm.DB) error {
	initSession(s)
	return nil
}

// Initialize some dynamic properties on the model, especially useful with GORM hooks
// for when a model is retrieved from the database
// TODO: add a method for mutating these properties to avoid this function
func initSession(s *Session) {
	s.ServerEvents = make(map[uint]chan ServerEvent)
	for _, p := range s.Players {
		s.ServerEvents[p.ID] = make(chan ServerEvent)
	}
	s.PlayerEvents = make(chan PlayerEvent)
}

func NewSession(db *gorm.DB, initialStage GameStage) (*Session, error) {
	return newSessionWithSeed(db, initialStage, time.Now().UnixNano)
}

func newSessionWithSeed(db *gorm.DB, initialStage GameStage, seed func() int64) (*Session, error) {
	var savedSession *Session
	for {
		rand.Seed(seed()) // TODO Use crypto.rand instead!
		savedSession = &Session{}
		session := Session{Code: alphaSessionCode(rand.Intn(SessionCodeMax))}
		result := db.FirstOrCreate(savedSession, &session)
		if result.Error != nil {
			return nil, result.Error
		} else if result.RowsAffected == 1 {
			break
		}
	}

	// Start the session in a go routine
	savedSession.CurrentStage = initialStage
	go savedSession.CurrentStage.Run(savedSession.PlayerEvents)

	return savedSession, nil
}

func (s *Session) AddPlayer(db *gorm.DB, p *Player) (chan ServerEvent, error) {
	s.Players = append(s.Players, *p)
	if result := db.Save(s); result.Error != nil {
		return nil, result.Error
	}

	c := make(chan ServerEvent, 100)
	var ctx context.Context // TODO get this
	s.PlayerEvents <- NewJoinEvent(ctx, p, c)
	return c, nil
}

func alphaSessionCode(code int) string {
	encoded := ""
	for len(encoded) < 4 {
		num := code % 26
		encoded = string(rune('A'+num)) + encoded
		code /= 26
	}
	return encoded
}
