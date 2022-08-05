package connect4

type Piece uint8

type Board [MaxRows][MaxColumns]Piece

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

const (
	Unclaimed Piece = iota
	Red
	Black
)

type GameResult uint8

const (
	GameWon GameResult = iota
	GameNotWon
)
