package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParticipant_Associations(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	kernel := &LambdaKernel{}
	session, _ := NewSession(db, kernel)
	pj, _ := NewParticipant(db, "Joe", session)
	pa, _ := NewParticipant(db, "Annie", session)

	savedSession := Session{}
	db.Preload("Participants").Find(&savedSession, session.ID)

	if len(savedSession.Participants) != 2 {
		t.Errorf("savedSession should have 2 participants; found: %d", len(savedSession.Participants))
	}
	for i, name := range []string{pj.Name, pa.Name} {
		if savedSession.Participants[i].Name != name {
			t.Errorf(`savedSession's participant at index %d should have been named "%s"`, i, name)
		}
	}
}

func TestNewParticipant_Channels(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	kernel := &LambdaKernel{}
	session, _ := NewSession(db, kernel)
	NewParticipant(db, "Joe", session)
	NewParticipant(db, "Annie", session)
	db.Preload("Participants").Find(&session, session.ID)

	assert.Len(t, session.Participants, 2, "there should have been 2 participants in the session")
	assert.Equal(t, len(session.ParticipantChannels), len(session.Participants), "there should be one channel per participant")
	for _, p := range session.Participants {
		assert.Containsf(t, session.ParticipantChannels, p.ID, "%s was not asigned a channel", p.Name)
		assert.NotNilf(t, session.ParticipantChannels[p.ID], "%s's channel is nil", p.Name)
	}
}
