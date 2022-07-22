package repositories

import (
	"testing"
	"time"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepo_NewSession(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionRepo{DB: db}
	kernel := &models.LambdaKernel{}
	session, _ := repo.NewSession(kernel)

	assert.Len(t, repo.LiveSessions, 1)
	assert.Equal(t, repo.LiveSessions[0].Code, session.Code)
}

func TestSessionRepo_SessionForCode(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionRepo{DB: db}
	session, _ := repo.NewSession(&models.LambdaKernel{})
	returnedSession, err := repo.SessionForCode(session.Code)

	assert.Equal(t, returnedSession.ID, session.ID)
	assert.Equal(t, returnedSession.Code, session.Code)
	assert.Nil(t, err)
}

func TestSessionRepo_SessionForCode_NotFound(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionRepo{DB: db}
	session, err := repo.SessionForCode("XXXX")

	assert.Error(t, err, `could not find session with code "XXXX"`)
	assert.Nil(t, session)
}

func TestSessionRepo_AddParticipantToSession(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionRepo{DB: db}
	kernel := &models.LambdaKernel{}
	session, _ := repo.NewSession(kernel)

	participant, _ := models.NewParticipant(db, "Steve", session)
	repo.AddParticipantToSession(participant, session.Code)

	assert.Equal(t, participant.Session.Code, session.Code, "Session association was not assigned on the participant")

	assert.Eventually(t, func() bool {
		return len(kernel.Events) == 1
	}, time.Second, 10*time.Millisecond)

	event := kernel.Events[0].(*models.JoinEvent)
	assert.Equal(t, event.Type(), models.JoinEventType, "The kernel should have received a Join event but didn't")
	assert.Equalf(t, event.Participant.Name, participant.Name, "The kernel received a Join event for the wrong participant")
}

func TestSessionRepo_AddParticipantToSession_ParticipantChannel(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionRepo{DB: db}
	kernel := models.NewWelcomeKernel()
	session, _ := repo.NewSession(kernel)

	names := []string{"Steve", "Angela"}
	participants := make([]*models.Participant, 0, len(names))
	for _, n := range names {
		p, _ := models.NewParticipant(db, n, session)
		repo.AddParticipantToSession(p, session.Code)
		participants = append(participants, p)
	}

	for _, p := range participants {
		assert.Equal(t, session.Code, p.Session.Code, "Session association was not assigned on the participant")
	}

	assert.Eventually(t, func() bool {
		return len(kernel.Events) == 2
	}, time.Second, 10*time.Millisecond)

	for i, event := range kernel.Events {
		event := event.(*models.JoinEvent)
		assert.Equal(t, models.JoinEventType, event.Type())
		assert.Equalf(t, names[i], event.Participant.Name, "Expected a join event for participant named %s", names[i])

		welcomeEvent := (<-session.ParticipantChannels[event.Participant.ID]).(*models.WelcomeEvent)
		assert.Equalf(t, welcomeEvent.Name, names[i], "Expected a welcome event for participant named %s", names[i])
	}
}
