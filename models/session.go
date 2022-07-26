package models

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const (
	SessionCodeLength = 4
	SessionCodeMax    = 456976 // 26^4
)

// TODO: move this to its own place
type GameKernel interface {
	Run(chan Event)
}

type Session struct {
	gorm.Model

	Code    string
	Players []Player

	Kernel         GameKernel          `gorm:"-:all"`
	PlayerChannels map[uint]chan Event `gorm:"-:all"`
	Events         chan Event          `gorm:"-:all"`
}

// Initialize some dynamic properties on the model, especially useful with GORM hooks
// for when a model is retrieved from the database
// TODO: add a method for mutating these properties to avoid this function
func initSession(s *Session) {
	s.PlayerChannels = make(map[uint]chan Event)
	for _, p := range s.Players {
		s.PlayerChannels[p.ID] = make(chan Event)
	}
	s.Events = make(chan Event)
}

func NewSession(db *gorm.DB, kernel GameKernel) (*Session, error) {
	return newSessionWithSeed(db, kernel, time.Now().UnixNano)
}

func newSessionWithSeed(db *gorm.DB, kernel GameKernel, seed func() int64) (*Session, error) {
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

	savedSession.Kernel = kernel
	go savedSession.Kernel.Run(savedSession.Events)

	return savedSession, nil
}

func (s *Session) AddPlayer(db *gorm.DB, p *Player) (chan Event, error) {
	s.Players = append(s.Players, *p)
	if result := db.Save(s); result.Error != nil {
		return nil, result.Error
	}

	c := make(chan Event, 100)
	s.Events <- &JoinEvent{
		Player:  p,
		Channel: c,
	}
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

func (s *Session) AfterCreate(tx *gorm.DB) error {
	initSession(s)
	return nil
}

func (s *Session) AfterFind(tx *gorm.DB) error {
	initSession(s)
	return nil
}
