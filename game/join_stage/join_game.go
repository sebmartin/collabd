package join_stage

import (
	"fmt"

	"github.com/sebmartin/collabd/models"
)

const (
	InitialPlayerArraySize = 10
)

type JoinGame struct {
	MinPlayers uint
	MaxPlayers uint
	StartGame  func([]*models.Player) models.StageRunner

	players []*models.Player
}

func (g *JoinGame) Run(playerEvents <-chan models.PlayerEvent) models.StageRunner {
	for event := range playerEvents {
		switch event := event.(type) {

		case *JoinEvent:
			handleJoin(event, g)

		case *StartEvent:
			next := handleStart(event, g)
			if next != nil {
				return next
			}
		}
	}

	panic("JoinGame stage's event loop ended without a StartEvent")
}

func handleJoin(event *JoinEvent, stage *JoinGame) {
	if len(stage.players) >= int(stage.MaxPlayers) {
		event.Sender().ServerEvents <- models.NewErrorEvent(
			fmt.Errorf("maximum player count of %d has already been reached", stage.MaxPlayers),
		)
		return
	}
	if stage.players == nil {
		stage.players = make([]*models.Player, 0, InitialPlayerArraySize)
	}
	stage.players = append(stage.players, event.Sender())
	event.Sender().ServerEvents <- NewDidJoinEvent(event.Sender())
}

func handleStart(event *StartEvent, stage *JoinGame) models.StageRunner {
	if len(stage.players) < int(stage.MinPlayers) {
		event.Sender().ServerEvents <- models.NewErrorEvent(
			fmt.Errorf("only %d player(s) have joined, a minimum of %d are required before the game can be started", len(stage.players), stage.MinPlayers),
		)
		return nil
	}

	for _, p := range stage.players {
		p.ServerEvents <- NewDidStartEvent(stage.players)
	}
	return stage.StartGame(stage.players)
}
