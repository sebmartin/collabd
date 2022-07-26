package models

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func predictableSeed() func() int64 {
	var last_seed int64
	return func() int64 {
		last_seed += 1
		return last_seed
	}
}

type textEvent string

func (e textEvent) Type() EventType {
	return "TEXT"
}

func (e textEvent) Sender() *Player {
	return nil
}

func TestNewSession(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	expected := "NBDX"
	session, _ := newSessionWithSeed(db, &LambdaState{}, predictableSeed())
	if session.Code != expected {
		t.Errorf(`NewSession() created session with code "%s"; expected "%s"`, session.Code, expected)
	}
	if session.ID == 0 {
		t.Error(`NewSession() returned session does not have a primary key`)
	}

	var count int64
	db.Model(&Session{}).Count(&count)
	if count != 1 {
		t.Errorf("Found %d total sessions, expected 1", count)
	}
}

func TestNewSession_CodeCollision(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	session1, _ := newSessionWithSeed(db, &LambdaState{}, predictableSeed())
	session2, _ := newSessionWithSeed(db, &LambdaState{}, predictableSeed())

	if session1.Code == session2.Code {
		t.Errorf(`Both sessions were created with code collision "%s"`, session1.Code)
	}

	var count int64
	db.Model(&Session{}).Count(&count)
	if count != 2 {
		t.Errorf("Found %d total sessions, expected 2", count)
	}
}

func Test_alphaSessionCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{name: "simple", code: 3, want: "AAAD"},
		{name: "min", code: 0, want: "AAAA"},
		{name: "max digit1", code: 26, want: "AABA"},
		{name: "max digit2", code: int(math.Pow(26, 2)) - 1, want: "AAZZ"},
		{name: "max digit3", code: int(math.Pow(26, 3)) - 1, want: "AZZZ"},
		{name: "max", code: int(math.Pow(26, 4)) - 1, want: "ZZZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alphaSessionCode(tt.code); got != tt.want {
				t.Errorf("alphaSessionCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionChannels(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	kernel := &LambdaState{}
	event := textEvent("Test event")
	session, _ := newSessionWithSeed(db, kernel, predictableSeed())
	session.PlayerEvents <- event

	assert.Len(t, kernel.Events, 1, "Expected exactly one event")

	receivedEvent := kernel.Events[0]
	assert.Equal(t, receivedEvent, event)
}
