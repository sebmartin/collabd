package connect4

import "fmt"

type Piece uint8

const (
	Unclaimed Piece = iota
	Red
	Black
)

const (
	MaxColumns uint = 7
	MaxRows    uint = 6
)

type Board [MaxRows][MaxColumns]Piece

func (board *Board) DropPiece(p Piece, slot uint) (int, error) {
	var targetHeight int
	if slot >= MaxColumns {
		return -1, fmt.Errorf("slot %d exceeds the slot maximum of %d", slot, MaxColumns)
	}
	for targetHeight = int(MaxRows) - 1; ; targetHeight-- {
		if targetHeight < 0 {
			return -1, fmt.Errorf("slot %d is full and cannot accept another piece", slot)
		}
		if board[targetHeight][slot] == Unclaimed {
			break
		}
	}
	board[targetHeight][slot] = p
	return int(targetHeight), nil
}

func (board *Board) AnalyzeMove(slot uint, row uint) GameResult {
	var droppedPiece Piece
	if board[row][slot] == Unclaimed {
		return GameNotWon
	} else {
		droppedPiece = board[row][slot]
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
				value := board[y][x]
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

func (b Board) String() string {
	var s string
	for y := 0; y < int(MaxRows); y++ {
		s += "["
		for x := 0; x < int(MaxColumns); x++ {
			switch b[y][x] {
			case Black:
				s += " B"
			case Red:
				s += " R"
			case Unclaimed:
				s += " -"
			}
		}
		s += " ]\n"
	}
	return s
}
