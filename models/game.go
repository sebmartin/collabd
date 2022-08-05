package models

type Game interface {
	Name() string
	InitialStage() *GameStage
}
