package models

import (
	"math/rand"
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

	Code    string
	Players []*Player `gorm:"-:all"` // TODO remove once session is no longer aware of players

	CurrentStage StageRunner               `gorm:"-:all"`
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
// TODO: maybe add a method for mutating these properties to avoid this function
func initSession(s *Session) {
	s.ServerEvents = make(map[uint]chan ServerEvent)
	for _, p := range s.Players {
		s.ServerEvents[p.ID] = make(chan ServerEvent, ChanBufferSize)
	}

	s.PlayerEvents = make(chan PlayerEvent, ChanBufferSize)
}

func NewSession(db *gorm.DB, initializer GameDescriber) (*Session, error) {
	return newSessionWithSeed(db, initializer, time.Now().UnixNano)
}

func newSessionWithSeed(db *gorm.DB, initializer GameDescriber, seed func() int64) (*Session, error) {
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

	// TODO - wrap this go routine in a lambda to manage the stage transitions
	// .. also, make that threadsafe
	go startSession(savedSession)

	return savedSession, nil
}

// This is the main game loop which is executed as a subroutine. It starts running
// the initial StageRunner and transitions to others as the runner processes events.
func startSession(session *Session) {
	currentStage := session.CurrentStage
	for {
		currentStage = currentStage.Run(session.PlayerEvents)

		// TODO: how does the game end?
		if currentStage == nil {
			break
		}
	}
}

func (s *Session) HandlePlayerEvent(event PlayerEvent) {
	s.PlayerEvents <- event
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
