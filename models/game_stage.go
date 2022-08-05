package models

type GameStage interface {
	Run(<-chan PlayerEvent) GameStage
}
