package connect4

import (
	"testing"
)

func TestConnect4_DropPiece(t *testing.T) {
	type args struct {
		p    Piece
		slot uint
	}
	tests := []struct {
		name       string
		board      Board
		args       args
		want       int
		wantErr    bool
		finalBoard Board
	}{
		{
			name: "blank, first column",
			args: args{
				p:    Red,
				slot: 0,
			},
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
			},
			finalBoard: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{R, X, X, X, X, X, X},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "blank, last column",
			args: args{
				p:    Red,
				slot: 6,
			},
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
			},
			finalBoard: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, R},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "stacked, same color",
			args: args{
				p:    Black,
				slot: 3,
			},
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{B, B, B, B, X, X, X},
			},
			finalBoard: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{B, B, B, B, X, X, X},
			},
			want:    4,
			wantErr: false,
		},
		{
			name: "stacked, different color",
			args: args{
				p:    Black,
				slot: 3,
			},
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{B, B, B, R, X, X, X},
			},
			finalBoard: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{B, B, B, R, X, X, X},
			},
			want:    4,
			wantErr: false,
		},
		{
			name: "column full",
			args: args{
				p:    Black,
				slot: 0,
			},
			board: Board{
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, B, B, R, X, X, X},
			},
			finalBoard: Board{
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{B, B, B, R, X, X, X},
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "slot too high",
			args: args{
				p:    Black,
				slot: MaxColumns,
			},
			board:      Board{},
			finalBoard: Board{},
			want:       -1,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Connect4{
				Board: tt.board,
			}
			got, err := g.DropPiece(tt.args.p, tt.args.slot)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect4.DropPiece() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Connect4.DropPiece() = %v, want %v", got, tt.want)
			}
			if g.Board != tt.finalBoard {
				t.Errorf("Final board did not match.\nExpected:\n%s\nGot:\n%s", tt.finalBoard, g.Board)
			}
		})
	}
}

const (
	X = Unclaimed
	R = Red
	B = Black
)

func TestConnect4_AnalyzeMove(t *testing.T) {
	type args struct {
		slot uint
		row  uint
	}
	tests := []struct {
		name  string
		board Board
		args  args
		want  GameResult
	}{
		{
			name: "horizontal on first row",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{B, B, B, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  5,
			},
			want: GameWon,
		},
		{
			name: "horizontal on first row; middle",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, B, B, B, B, X},
			},
			args: args{
				slot: 3,
				row:  5,
			},
			want: GameWon,
		},
		{
			name: "horizontal on first row; right edge",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, B, B, B},
			},
			args: args{
				slot: 3,
				row:  5,
			},
			want: GameWon,
		},
		{
			name: "horizontal on first row; intercepted other color",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{B, R, B, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  5,
			},
			want: GameNotWon,
		},
		{
			name: "horizontal on first row; intercepted blank",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{B, X, B, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  5,
			},
			want: GameNotWon,
		},
		{
			name: "diagonal \\",
			board: Board{
				{X, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{R, B, X, X, X, X, X},
				{B, R, B, X, X, X, X},
				{R, R, B, B, X, X, X},
				{B, B, R, R, X, X, X},
			},
			args: args{
				slot: 1,
				row:  2,
			},
			want: GameWon,
		},
		{
			name: "diagonal \\; intercepted with other color",
			board: Board{
				{X, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{R, B, X, X, X, X, X},
				{B, R, R, X, X, X, X},
				{R, R, B, B, X, X, X},
				{B, B, R, R, X, X, X},
			},
			args: args{
				slot: 1,
				row:  2,
			},
			want: GameNotWon,
		},
		{
			name: "diagonal \\; intercepted with blank",
			board: Board{
				{X, X, X, X, X, X, X},
				{B, X, X, X, X, X, X},
				{R, B, X, X, X, X, X},
				{B, R, X, X, X, X, X},
				{R, R, B, B, X, X, X},
				{B, B, R, R, X, X, X},
			},
			args: args{
				slot: 1,
				row:  2,
			},
			want: GameNotWon,
		},
		{
			name: "diagonal \\ niddle",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, B, X, X, X, X},
				{X, X, R, B, X, X, X},
				{X, X, B, R, B, X, X},
				{X, X, R, R, B, B, X},
				{X, X, B, B, R, R, X},
			},
			args: args{
				slot: 2,
				row:  1,
			},
			want: GameWon,
		},
		{
			name: "diagonal \\ right",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, B, X, X},
				{X, X, X, B, R, B, X},
				{X, X, X, R, R, B, B},
				{X, X, X, B, B, R, R},
			},
			args: args{
				slot: 6,
				row:  4,
			},
			want: GameWon,
		},
		{
			name: "diagonal /",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, B, R, X, X, X},
				{X, B, R, B, X, X, X},
				{B, B, R, R, X, X, X},
				{R, R, B, B, X, X, X},
			},
			args: args{
				slot: 0,
				row:  4,
			},
			want: GameWon,
		},
		{
			name: "diagonal /; intercepted with other color",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, B, R, X, X, X},
				{X, R, R, B, X, X, X},
				{B, B, R, R, X, X, X},
				{R, R, B, B, X, X, X},
			},
			args: args{
				slot: 2,
				row:  2,
			},
			want: GameNotWon,
		},
		{
			name: "diagonal /; intercepted with blank",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, X, X, X},
				{X, B, R, B, X, X, X},
				{B, B, R, R, X, X, X},
				{R, R, B, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  1,
			},
			want: GameNotWon,
		},
		{
			name: "diagonal / niddle",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, B, X},
				{X, X, X, X, B, R, X},
				{X, X, X, B, R, B, X},
				{X, X, B, B, R, R, X},
				{X, X, R, R, B, B, X},
			},
			args: args{
				slot: 2,
				row:  4,
			},
			want: GameWon,
		},
		{
			name: "diagonal / right",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, B},
				{X, X, X, X, X, B, R},
				{X, X, X, X, B, R, B},
				{X, X, X, B, B, R, R},
				{X, X, X, R, R, B, B},
			},
			args: args{
				slot: 6,
				row:  1,
			},
			want: GameWon,
		},
		{
			name: "vertical",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  3,
			},
			want: GameWon,
		},
		{
			name: "vertical top",
			board: Board{
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, X, X, X},
				{X, X, X, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  3,
			},
			want: GameWon,
		},
		{
			name: "vertical; incomplete",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  3,
			},
			want: GameNotWon,
		},
		{
			name: "vertical; intercepted 1",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, X, X, X},
			},
			args: args{
				slot: 3,
				row:  3,
			},
			want: GameNotWon,
		},
		{
			name: "vertical; intercepted 2",
			board: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, X, X, X},
				{X, X, X, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  3,
			},
			want: GameNotWon,
		},
		{
			name: "vertical; intercepted 3",
			board: Board{
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, R, X, X, X},
				{X, X, X, B, X, X, X},
				{X, X, X, B, X, X, X},
			},
			args: args{
				slot: 3,
				row:  1,
			},
			want: GameNotWon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Connect4{
				Board: tt.board,
			}
			if got := g.AnalyzeMove(tt.args.slot, tt.args.row); got != tt.want {
				t.Errorf("Connect4.AnalyzeMove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoard_String(t *testing.T) {
	tests := []struct {
		name string
		b    Board
		want string
	}{
		{
			name: "complex",
			b: Board{
				{X, X, X, X, X, X, X},
				{X, X, X, X, X, X, B},
				{X, X, X, X, X, B, R},
				{X, X, X, X, B, R, B},
				{X, X, X, B, B, R, R},
				{X, X, X, R, R, B, B},
			},
			want: `[ - - - - - - - ]
[ - - - - - - B ]
[ - - - - - B R ]
[ - - - - B R B ]
[ - - - B B R R ]
[ - - - R R B B ]
`,
		},
		{
			name: "checkers",
			b: Board{
				{X, X, X, X, X, X, X},
				{B, R, B, R, B, R, B},
				{R, B, R, B, R, B, R},
				{B, R, B, R, B, R, B},
				{R, B, R, B, R, B, R},
				{B, R, B, R, B, R, B},
			},
			want: `[ - - - - - - - ]
[ B R B R B R B ]
[ R B R B R B R ]
[ B R B R B R B ]
[ R B R B R B R ]
[ B R B R B R B ]
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Board.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
