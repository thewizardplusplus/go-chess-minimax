package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/generators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNegamaxSearcher(test *testing.T) {
	type args struct {
		maximalDeep int
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
				pieces:      []models.Piece{},
				color:       models.White,
			},
			wantMove: ScoredMove{},
			wantErr:  nil,
		},
	} {
		generator :=
			generators.NewDefaultMoveGenerator(
				models.MoveGenerator{},
			)
		terminator :=
			terminators.NewDeepTerminator(
				data.args.maximalDeep,
			)
		board := models.NewBoard(
			models.Size{8, 8},
			data.args.pieces,
		)
		searcher := NewNegamaxSearcher(
			generator,
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
