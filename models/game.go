package models

type Game struct {
	name         string
	initialStage StageRunner
}

func (g Game) Name() string {
	return g.name
}

func (g Game) InitialStage() StageRunner {
	return g.initialStage
}

func NewGame(name string, stage StageRunner) *Game {
	return &Game{
		name:         name,
		initialStage: stage,
	}
}

type GameDescriber interface {
	Name() string
	InitialStage() StageRunner
}
