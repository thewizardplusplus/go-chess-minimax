package chessminimax

import (
	"reflect"
	"testing"

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

	for _, data := range []data{} {
		generator := NewDefaultMoveGenerator(
			models.MoveGenerator{},
		)
		terminator := NewDeepTerminator(
			data.args.maximalDeep,
		)
		board := models.NewBoard(
			models.Size{8, 8},
			data.args.pieces,
		)
		searcher := NewNegamaxSearcher(
			generator,
			terminator,
			MaterialEvaluator{},
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
			test.Fail()
		}
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Fail()
		}
	}
}
