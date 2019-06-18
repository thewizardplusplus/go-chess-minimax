package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func TestNegamaxSearcher(test *testing.T) {
	type args struct {
		maximalDeep int
		size        models.Size
		pieces      []models.Piece
		color       models.Color
	}
	type data struct {
		args     args
		wantMove ScoredMove
		wantErr  error
	}

	for _, data := range []data{
		data{
			args: args{
				maximalDeep: 0,
				size:        models.Size{8, 8},
				pieces: []models.Piece{
					pieces.NewKing(
						models.Black,
						models.Position{0, 0},
					),
					pieces.NewKing(
						models.White,
						models.Position{7, 7},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 0},
					),
				},
				color: models.White,
			},
			wantMove: ScoredMove{},
			wantErr:  models.ErrKingCapture,
		},
		data{
			args: args{
				maximalDeep: 0,
				size:        models.Size{8, 8},
				pieces: []models.Piece{
					pieces.NewKing(
						models.Black,
						models.Position{0, 0},
					),
					pieces.NewKing(
						models.White,
						models.Position{7, 7},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 0},
					),
				},
				color: models.Black,
			},
			wantMove: ScoredMove{Score: -5},
			wantErr:  nil,
		},
		data{
			args: args{
				maximalDeep: 1,
				size:        models.Size{8, 8},
				pieces: []models.Piece{
					pieces.NewKing(
						models.Black,
						models.Position{0, 0},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{0, 1},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{1, 1},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{1, 0},
					),
					pieces.NewKing(
						models.White,
						models.Position{7, 7},
					),
				},
				color: models.Black,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
		data{
			args: args{
				maximalDeep: 1,
				size:        models.Size{8, 8},
				pieces: []models.Piece{
					pieces.NewKing(
						models.Black,
						models.Position{0, 1},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{0, 2},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{1, 2},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{1, 1},
					),
					pieces.NewKing(
						models.White,
						models.Position{7, 7},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 0},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 1},
					),
				},
				color: models.Black,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
		data{
			args: args{
				maximalDeep: 1,
				size:        models.Size{8, 8},
				pieces: []models.Piece{
					pieces.NewKing(
						models.Black,
						models.Position{0, 1},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{0, 2},
					),
					pieces.NewPawn(
						models.Black,
						models.Position{1, 2},
					),
					pieces.NewKing(
						models.White,
						models.Position{7, 7},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 0},
					),
					pieces.NewRook(
						models.White,
						models.Position{7, 1},
					),
					pieces.NewBishop(
						models.White,
						models.Position{6, 7},
					),
				},
				color: models.Black,
			},
			wantMove: ScoredMove{
				Score: evaluateCheckmate(0),
			},
			wantErr: ErrCheckmate,
		},
	} {
		terminator :=
			terminators.NewDeepTerminator(
				data.args.maximalDeep,
			)
		board := models.NewBoard(
			data.args.size,
			data.args.pieces,
		)
		searcher := NewNegamaxSearcher(
			models.MoveGenerator{},
			terminator,
			evaluators.MaterialEvaluator{},
		)
		gotMove, gotErr := searcher.SearchMove(
			board,
			data.args.color,
			0,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Log(gotMove)
			test.Log(data.wantMove)
			test.Fail()
		}
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Log(gotErr)
			test.Log(data.wantErr)
			test.Fail()
		}
	}
}
