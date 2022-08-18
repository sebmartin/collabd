package models

type Game struct {
	name         string
	initialStage StageRunner
}

func NewGame(name string, stage StageRunner) *Game {
	return &Game{
		name:         name,
		initialStage: stage,
	}
}

func (g Game) Name() string {
	return g.name
}

func (g Game) InitialStage() StageRunner {
	return g.initialStage
}

// TODO this is exploration
type GameInitializer interface {
	Name() string
	InitialStage() StageRunner
}
