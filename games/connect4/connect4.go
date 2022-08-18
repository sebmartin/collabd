package connect4

import (
	"context"
	"fmt"

	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/models"
)

func init() {
	game.Register("Connect4", func(ctx context.Context) (models.GameInitializer, error) {
		return NewGame(), nil
	})
}

const (
	MaxColumns uint = 7
	MaxRows    uint = 6
)

type Connect4 struct {
	models.Game

	Board Board
}

func NewGame() Connect4 {
	return Connect4{
		Game: *models.NewGame(
			"Connect 4",
			nil, // TODO: set initialStage
		),
		Board: [6][7]Piece{},
	}
}

func (g *Connect4) DropPiece(p Piece, slot uint) (int, error) {
	var targetHeight int
	if slot >= MaxColumns {
		return -1, fmt.Errorf("slot %d exceeds the slot maximum of %d", slot, MaxColumns)
	}
	for targetHeight = int(MaxRows) - 1; ; targetHeight-- {
		if targetHeight < 0 {
			return -1, fmt.Errorf("slot %d is full and cannot accept another piece", slot)
		}
		if g.Board[targetHeight][slot] == Unclaimed {
			break
		}
	}
	g.Board[targetHeight][slot] = p
	return int(targetHeight), nil
}

func (g *Connect4) AnalyzeMove(slot uint, row uint) GameResult {
	var droppedPiece Piece
	if g.Board[row][slot] == Unclaimed {
		return GameNotWon
	} else {
		droppedPiece = g.Board[row][slot]
	}

	vectors := []func(int) [2]int{
		func(i int) [2]int { return [2]int{0, i} },  // Vertical
		func(i int) [2]int { return [2]int{i, 0} },  // Horizontal
		func(i int) [2]int { return [2]int{i, i} },  // Diagonal /
		func(i int) [2]int { return [2]int{-i, i} }, // Diagonal \
	}

	for _, f := range vectors {
		count := 1
		for _, dir := range []int{1, -1} {
			for i := dir; i < 4 && i > -4; i += dir {

				v := f(i)
				x, y := int(slot)+v[0], int(row)+v[1]
				if x < 0 || y < 0 || x >= int(MaxColumns) || y >= int(MaxRows) {
					break
				}
				value := g.Board[y][x]
				if value == droppedPiece {
					count += 1
				} else {
					break
				}

				if count >= 4 {
					return GameWon
				}
			}
		}
	}
	return GameNotWon
}
