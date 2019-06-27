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
		boardInFEN  string
		color       models.Color
		maximalDeep int
	}
	type data struct {
		args     args
		wantMove ScoredMove
		wantErr  error
	}

	for _, data := range []data{
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/8/k6R",
				color:       models.White,
				maximalDeep: 0,
			},
			wantMove: ScoredMove{},
			wantErr:  models.ErrKingCapture,
		},
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/8/k6R",
				color:       models.Black,
				maximalDeep: 0,
			},
			wantMove: ScoredMove{Score: -5},
			wantErr:  nil,
		},
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/pp6/kp6",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
		data{
			args: args{
				boardInFEN: "7K/8/8/8" +
					"/8/pp6/kp5R/7R",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
		data{
			args: args{
				boardInFEN: "6BK/8/8/8" +
					"/8/pp6/k6R/7R",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: ScoredMove{
				Score: evaluateCheckmate(0),
			},
			wantErr: ErrCheckmate,
		},
	} {
		storage, err := models.ParseBoard(
			data.args.boardInFEN,
			pieces.NewPiece,
		)
		if err != nil {
			test.Fail()
			continue
		}

		generator := models.MoveGenerator{}
		terminator :=
			terminators.NewDeepTerminator(
				data.args.maximalDeep,
			)
		evaluator :=
			evaluators.MaterialEvaluator{}
		searcher := NewNegamaxSearcher(
			generator,
			terminator,
			evaluator,
		)
		gotMove, gotErr := searcher.SearchMove(
			storage,
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
