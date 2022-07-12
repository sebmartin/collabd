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

	Code string
}

func NewSession(db *gorm.DB) *Session {
	return newSessionWithSeed(db, time.Now().UnixNano)
}

func newSessionWithSeed(db *gorm.DB, seed func() int64) *Session {
	for {
		rand.Seed(seed())
		session := &Session{Code: alphaSessionCode(rand.Intn(SessionCodeMax))}
		if result := db.FirstOrCreate(&Session{}, session); result.RowsAffected == 1 {
			return session
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
