package services

import (
	"testing"
	"time"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepo_NewSession(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionService{DB: db}
	kernel := &models.LambdaState{}
	session, _ := repo.NewSession(kernel)

	assert.Len(t, repo.LiveSessions, 1)
	assert.Equal(t, repo.LiveSessions[0].Code, session.Code)
}

func TestSessionRepo_SessionForCode(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionService{DB: db}
	session, _ := repo.NewSession(&models.LambdaState{})
	returnedSession, err := repo.SessionForCode(session.Code)

	assert.Equal(t, returnedSession.ID, session.ID)
	assert.Equal(t, returnedSession.Code, session.Code)
	assert.Nil(t, err)
}

func TestSessionRepo_SessionForCode_NotFound(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionService{DB: db}
	session, err := repo.SessionForCode("XXXX")

	assert.Error(t, err, `could not find session with code "XXXX"`)
	assert.Nil(t, session)
}

func TestSessionRepo_JoinSession(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionService{DB: db}
	kernel := &models.LambdaState{}
	session, _ := repo.NewSession(kernel)

	player, _ := models.NewPlayer(db, "Steve", session)
	repo.JoinSession(player, session.Code)

	assert.Equal(t, player.Session.Code, session.Code, "Session association was not assigned on the player")

	assert.Eventually(t, func() bool {
		return len(kernel.Events) == 1
	}, time.Second, 10*time.Millisecond)

	event := kernel.Events[0].(*models.JoinEvent)
	assert.Equal(t, event.Type(), models.JoinEventType, "The kernel should have received a Join event but didn't")
	assert.Equalf(t, event.Sender().Name, player.Name, "The kernel received a Join event for the wrong player")
}

func TestSessionRepo_JoinSession_PlayerChannel(t *testing.T) {
	db, cleanup := models.ConnectWithTestDB()
	defer cleanup()

	repo := SessionService{DB: db}
	kernel := models.NewWelcomeKernel()
	session, _ := repo.NewSession(kernel)

	names := []string{"Steve", "Angela"}
	players := make([]*models.Player, 0, len(names))
	for _, n := range names {
		p, _ := models.NewPlayer(db, n, session)
		repo.JoinSession(p, session.Code)
		players = append(players, p)
	}

	for _, p := range players {
		assert.Equal(t, session.Code, p.Session.Code, "Session association was not assigned on the player")
	}

	assert.Eventually(t, func() bool {
		return len(kernel.Events) == 2
	}, time.Second, 10*time.Millisecond)

	for i, event := range kernel.Events {
		event := event.(*models.JoinEvent)
		assert.Equal(t, models.JoinEventType, event.Type())
		assert.Equalf(t, names[i], event.Sender().Name, "Expected a join event for player named %s", names[i])

		welcomeEvent := (<-session.ServerEvents[event.Sender().ID]).(*models.WelcomeEvent)
		assert.Equalf(t, welcomeEvent.Name, names[i], "Expected a welcome event for player named %s", names[i])
	}
}
