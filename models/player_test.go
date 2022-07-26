package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayer_Associations(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	kernel := &LambdaState{}
	session, _ := NewSession(db, kernel)
	pj, _ := NewPlayer(db, "Joe", session)
	pa, _ := NewPlayer(db, "Annie", session)

	savedSession := Session{}
	db.Preload("Players").Find(&savedSession, session.ID)

	if len(savedSession.Players) != 2 {
		t.Errorf("savedSession should have 2 players; found: %d", len(savedSession.Players))
	}
	for i, name := range []string{pj.Name, pa.Name} {
		if savedSession.Players[i].Name != name {
			t.Errorf(`savedSession's player at index %d should have been named "%s"`, i, name)
		}
	}
}

func TestNewPlayer_Channels(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	kernel := &LambdaState{}
	session, _ := NewSession(db, kernel)
	NewPlayer(db, "Joe", session)
	NewPlayer(db, "Annie", session)
	db.Preload("Players").Find(&session, session.ID)

	assert.Len(t, session.Players, 2, "there should have been 2 players in the session")
	assert.Equal(t, len(session.ServerEvents), len(session.Players), "there should be one channel per player")
	for _, p := range session.Players {
		assert.Containsf(t, session.ServerEvents, p.ID, "%s was not asigned a channel", p.Name)
		assert.NotNilf(t, session.ServerEvents[p.ID], "%s's channel is nil", p.Name)
	}
}
