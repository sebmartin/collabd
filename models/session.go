package models

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"gorm.io/gorm"
)

const (
	SessionCodeLength = 4
	SessionCodeMax    = 456976 // 26^4
	ChanBufferSize    = 100
)

type contextKey string

const SessionKey = contextKey("session")

type Session struct {
	gorm.Model

	Code      string
	Players   []*Player    `gorm:"-:all"`
	PlayersMu sync.RWMutex `gorm:"-:all"` // TODO: this sucks, if the game routine modifies this as part of a stage then we don't have to worry about synchronization

	CurrentStage   StageRunner               `gorm:"-:all"`
	ServerEvents   map[uint]chan ServerEvent `gorm:"-:all"`
	ServerEventsMu sync.RWMutex              `gorm:"-:all"`
	PlayerEvents   chan PlayerEventEnvelope  `gorm:"-:all"`
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
		s.ServerEvents[p.ID] = make(chan ServerEvent, ChanBufferSize)
	}

	s.PlayerEvents = make(chan PlayerEventEnvelope, ChanBufferSize)
}

func NewSession(db *gorm.DB, initializer GameInitializer) (*Session, error) {
	return newSessionWithSeed(db, initializer, time.Now().UnixNano)
}

func newSessionWithSeed(db *gorm.DB, initializer GameInitializer, seed func() int64) (*Session, error) {
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
	savedSession.CurrentStage = initializer.InitialStage()
	go savedSession.CurrentStage.Run(savedSession.PlayerEvents)

	return savedSession, nil
}

// TODO get DB from context
func (s *Session) AddPlayer(ctx context.Context, db *gorm.DB, p *Player) error {
	s.Players = append(s.Players, p)
	if result := db.Save(s); result.Error != nil {
		return result.Error
	}

	s.ServerEvents[p.ID] = p.ServerEvents
	s.HandlePlayerEvent(ctx, NewJoinEvent(ctx, p, p.ServerEvents))
	// s.PlayerEvents <- NewJoinEvent(ctx, p, p.ServerEvents)
	return nil
}

func (s *Session) HandlePlayerEvent(ctx context.Context, event PlayerEvent) {
	s.PlayerEvents <- PlayerEventEnvelope{
		PlayerEvent: event,
		Session:     *s,
		Context:     context.WithValue(ctx, SessionKey, s),
	}
}

func (s *Session) SendServerEvent(playerID uint, event ServerEvent) {
	panic("// TODO")
}

func (s *Session) BoardcastServerEvent(event ServerEvent) {
	for _, p := range s.Players {
		s.SendServerEvent(p.ID, event)
	}
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
