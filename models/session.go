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

type Session struct {
	gorm.Model

	Code         string
	Participants []Participant
}

func NewSession(db *gorm.DB) (*Session, error) {
	return newSessionWithSeed(db, time.Now().UnixNano)
}

// improve error handling, return an error tuple
func newSessionWithSeed(db *gorm.DB, seed func() int64) (*Session, error) {
	for {
		rand.Seed(seed())
		session := Session{Code: alphaSessionCode(rand.Intn(SessionCodeMax))}
		savedSession := &Session{}
		result := db.FirstOrCreate(savedSession, &session)
		if result.Error != nil {
			return nil, result.Error
		} else if result.RowsAffected == 1 {
			return savedSession, nil
		}
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
