package chessminimax

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNegamaxSearcher(test *testing.T) {
	type data struct {
		maximalDeep int
		pieces      []models.Piece
		color       models.Color
		wantMove    ScoredMove
		wantErr     error
	}

	for _, data := range []data{} {
		generator := NewDefaultMoveGenerator(
			models.MoveGenerator{},
		)
		terminator :=
			NewDeepTerminator(data.maximalDeep)
		board := models.NewBoard(
			models.Size{8, 8},
			data.pieces,
		)
		searcher := NewNegamaxSearcher(
			generator,
			terminator,
			MaterialEvaluator{},
		)
		gotMove, gotErr := searcher.SearchMove(
			board,
			data.color,
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
