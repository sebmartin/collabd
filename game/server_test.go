package game

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
)

const testGameName = "__test_game__"

func init() {
	registerTestGame()
}

func registerTestGame() {
	Register(testGameName, func(ctx context.Context) (models.GameInitializer, error) {
		return TestGame{
			Game: *models.NewGame(
				"Test Game",
				&TestStage{},
			),
		}, nil
	})
}

func newServer(t *testing.T) (*Server, func()) {
	dbtmpdir, _ := os.MkdirTemp("", "collabd_tests_")
	dbtmppath := path.Join(dbtmpdir, "_tests.sqlite")
	server, err := NewServer("sqlite", dbtmppath)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %s", err)
	}
	return server, func() {
		sqldb, err := server.db.DB()
		if err != nil {
			sqldb.Close()
		}
		os.Remove(dbtmppath)
		os.Remove(dbtmpdir)
	}
}

func newServerSession(t *testing.T) (*Server, *models.Session, func()) {
	server, cleanup := newServer(t)

	// stage := &models.LambdaStage{}
	session, _ := server.NewSession(context.Background(), testGameName)
	return server, session, cleanup
}

func TestServer_NewSession(t *testing.T) {
	server, session, cleanup := newServerSession(t)
	defer cleanup()

	assert.Len(t, server.sessions, 1)
	assert.Equal(t, session, server.sessions[0])

	fetched, err := server.SessionForCode(session.Code)
	assert.Nilf(t, err, "Session could not be retrieved by code: %s", err)
	assert.Equal(t, session.Code, fetched.Code)
	assert.Equal(t, session, fetched)
}

func TestServer_NewSession_UnknownGame(t *testing.T) {
	server, cleanup := newServer(t)
	defer cleanup()

	_, err := server.NewSession(context.Background(), "UNKNOWN_GAME_NAME")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, `failed to create session, unknown game: UNKNOWN_GAME_NAME`)
}

func TestServer_SessionForCode_UnknownCode(t *testing.T) {
	server, cleanup := newServer(t)
	defer cleanup()

	_, err := server.SessionForCode("ABCD")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, `could not find session with code "ABCD"`)
}

func TestServer_JoinSession(t *testing.T) {
	server, session, cleanup := newServerSession(t)
	defer cleanup()

	player, _ := models.NewPlayer(server.db, "Steve")
	joinedSession, err := server.JoinSession(context.Background(), player, session.Code)
	assert.Nil(t, err)
	assert.Equal(t, session, joinedSession)
	assert.Len(t, joinedSession.Players, 1)
	assert.Len(t, session.Players, 1)
	assert.Equal(t, joinedSession.Players[0], player)
}

func TestServer_JoinSession_ServerEventsChannel(t *testing.T) {
	// TODO: this is testing code in session so it might make more sense to test it in that context
	server, session, cleanup := newServerSession(t)
	defer cleanup()

	player, _ := models.NewPlayer(server.db, "Steve")
	joinedSession, _ := server.JoinSession(context.Background(), player, session.Code)

	assert.Contains(t, joinedSession.ServerEvents, player.ID, "No player event channel found in session for player")

	event := &models.WelcomeEvent{
		Name: player.Name,
	}

	// Send dummy event
	select {
	case joinedSession.ServerEvents[player.ID] <- event:
		break
	default:
		assert.FailNow(t, "Player's server events channel is unbuffered (blocked)")
	}

	// Receive the event
	select {
	case rcvEvent := <-player.ServerEvents:
		assert.Equal(t, event, rcvEvent)
	default:
		assert.FailNow(t, "Server event was not received")
	}
}

func TestServer_JoinSession_PlayerEventsChannels(t *testing.T) {
	// server, session, cleanup := newServerSession(t)
	// defer cleanup()

	// player, _ := models.NewPlayer(server.db, "Steve")
	// joinedSession, _ := server.JoinSession(context.Background(), player, session.Code)

	// 	TODO make sure player channels work
}

func TestServer_JoinSession_UnknownCode(t *testing.T) {
	server, cleanup := newServer(t)
	defer cleanup()

	player, _ := models.NewPlayer(server.db, "Steve")
	joinedSession, err := server.JoinSession(context.Background(), player, "ABCD")
	assert.Nil(t, joinedSession)
	assert.ErrorContains(t, err, `could not find session with code "ABCD"`)
}

// - Fixtures

type TestGame struct {
	models.Game
}

type TestStage struct {
	recvPlayerEvents []models.PlayerEvent // Not thread safe, ok for now
}

func (s *TestStage) Run(playerEvents <-chan models.PlayerEventEnvelope) models.StageRunner {
	if s.recvPlayerEvents == nil {
		s.recvPlayerEvents = make([]models.PlayerEvent, 10)
	}

	for {
		event, ok := <-playerEvents
		if !ok {
			return nil
		}
		s.recvPlayerEvents = append(s.recvPlayerEvents, event.PlayerEvent)

		if event.Type() == models.EventType("ECHO") {
			event.Session.SendServerEvent(event.Sender().ID, event.PlayerEvent)
		}
	}
}

type EchoEvent struct {
	Name    string
	Message string
}

func (e *EchoEvent) Type() models.EventType {
	return models.EventType("ECHO")
}
