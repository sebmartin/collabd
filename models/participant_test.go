package models

import (
	"testing"
)

func TestNewParticipant(t *testing.T) {
	db, cleanup := ConnectWithTestDB()
	defer cleanup()

	session, _ := NewSession(db)
	pj, _ := NewParticipant(db, "Joe", session.Code)
	pa, _ := NewParticipant(db, "Annie", session.Code)

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
