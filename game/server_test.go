package game

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/sebmartin/collabd/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const testGameName = "__test_game__"

func init() {
	registerTestGame()
}

func registerTestGame() {
	Register(testGameName, func(ctx context.Context) (models.GameDescriber, error) {
		return testGame{
			Game: *models.NewGame(
				"Test Game",
				&testStage{},
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

	session, _ := server.NewSession(context.Background(), testGameName)
	return server, session, cleanup
}

func newPlayer(db *gorm.DB, name string) *models.Player {
	player, _ := models.NewPlayer(db, name)
	return player
}

func TestServer_NewSession_SessionForCode(t *testing.T) {
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

func TestServer_HandlePlayerEvent(t *testing.T) {
	server, session, cleanup := newServerSession(t)
	defer cleanup()

	player, _ := models.NewPlayer(server.db, "Steve")
	event := newEchoEvent(context.Background(), "Well hello there!", player)
	result := server.HandlePlayerEvent(session.Code, event)

	assert.Nil(t, result)

	select {
	case serverEvent := <-player.ServerEvents:
		assert.IsType(t, &echoEchoEvent{}, serverEvent)
		assert.Equal(t, event, serverEvent.(*echoEchoEvent).OriginalEvent)
	case <-time.After(500 * time.Millisecond):
		require.Fail(t, "Timeout", "Did not receive expected server event before timeout")
	}
}

func TestServer_HandlePlayerEvent_UnknownCode(t *testing.T) {
	server, _, cleanup := newServerSession(t)
	defer cleanup()

	player := newPlayer(server.db, "Steve")
	event := newEchoEvent(context.Background(), "Well hello there!", player)
	result := server.HandlePlayerEvent("XXXX", event)

	assert.ErrorContains(t, result, `could not find session with code "XXXX"`)
}

func TestBroadcast(t *testing.T) {
	server, _ := newServer(t)
	players := []*models.Player{
		newPlayer(server.db, "Alice"),
		newPlayer(server.db, "John"),
		newPlayer(server.db, "Sophie"),
	}

	echoEvent := newEchoEvent(context.Background(), "hello", players[0])
	Broadcast(players, newEchoEchoEvent(echoEvent))

	for _, p := range players {
		select {
		case event := <-p.ServerEvents:
			assert.IsType(t, &echoEchoEvent{}, event)
		default:
			assert.Fail(t, "Did not receive the event", "Player: %s", p.Name)
		}
	}
}

// - Fixtures

type testGame struct {
	models.Game
}

type testStage struct{}

func (s *testStage) Run(playerEvents <-chan models.PlayerEvent) models.StageRunner {
	for {
		event := <-playerEvents
		switch event := event.(type) {
		case *echoEvent:
			event.Sender().ServerEvents <- newEchoEchoEvent(event)
		}
	}
}

type echoEvent struct {
	models.PlayerEvent

	Message string
}

func newEchoEvent(ctx context.Context, message string, sender *models.Player) *echoEvent {
	return &echoEvent{
		PlayerEvent: models.NewPlayerEvent(ctx, "ECHO", sender),
		Message:     message,
	}
}

type echoEchoEvent struct {
	models.ServerEvent

	OriginalEvent *echoEvent
}

func newEchoEchoEvent(event *echoEvent) *echoEchoEvent {
	return &echoEchoEvent{
		ServerEvent:   models.NewServerEvent("ECHO_ECHO"),
		OriginalEvent: event,
	}
}
