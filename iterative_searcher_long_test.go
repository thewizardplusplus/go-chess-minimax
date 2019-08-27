// +build long

package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func TestIterativeSearcher(
	test *testing.T,
) {
	type args struct {
		boardInFEN  string
		color       models.Color
		maximalDeep int
	}
	type data struct {
		args     args
		wantMove moves.ScoredMove
		wantErr  error
	}

	for index, data := range []data{
		// king capture
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/8/k6R",
				color:       models.White,
				maximalDeep: 0,
			},
			wantMove: moves.ScoredMove{},
			wantErr:  models.ErrKingCapture,
		},
		// termination
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/8/k6R",
				color:       models.Black,
				maximalDeep: 0,
			},
			wantMove: moves.ScoredMove{Score: -5},
			wantErr:  nil,
		},
		// draw without checks
		data{
			args: args{
				boardInFEN:  "7K/8/8/8/8/8/pp6/kp6",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: moves.ScoredMove{},
			wantErr:  ErrDraw,
		},
		// draw with checks on a first ply
		data{
			args: args{
				boardInFEN: "7K/8/8/8" +
					"/8/pp6/kp5R/7R",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: moves.ScoredMove{},
			wantErr:  ErrDraw,
		},
		// draw with checks on a third ply
		data{
			args: args{
				boardInFEN: "7K/6P1/8/2q5" +
					"/8/8/b7/kb2B3",
				color:       models.White,
				maximalDeep: 3,
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start:  models.Position{4, 0},
					Finish: models.Position{2, 2},
				},
				Score: 0,
			},
			wantErr: nil,
		},
		// checkmate on a first ply
		data{
			args: args{
				boardInFEN: "6BK/8/8/8" +
					"/8/pp6/k6R/7R",
				color:       models.Black,
				maximalDeep: 1,
			},
			wantMove: moves.ScoredMove{
				Score: evaluateCheckmate(0),
			},
			wantErr: ErrCheckmate,
		},
		// checkmate on a second ply
		data{
			args: args{
				boardInFEN: "6K1/8/7q/6p1" +
					"/8/2B5/pp4PQ/k7",
				color:       models.White,
				maximalDeep: 2,
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start:  models.Position{7, 1},
					Finish: models.Position{6, 0},
				},
				Score: -evaluateCheckmate(1),
			},
			wantErr: nil,
		},
		// single legal move
		data{
			args: args{
				boardInFEN:  "7K/8/7q/8/8/8/8/k7",
				color:       models.White,
				maximalDeep: 1,
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start:  models.Position{7, 7},
					Finish: models.Position{6, 7},
				},
				Score: -9,
			},
			wantErr: nil,
		},
		// single profitable move on a first ply
		data{
			args: args{
				boardInFEN:  "7K/8/7q/8/8/8/7Q/k7",
				color:       models.White,
				maximalDeep: 1,
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start:  models.Position{7, 1},
					Finish: models.Position{7, 5},
				},
				Score: 9,
			},
			wantErr: nil,
		},
		// single profitable move on a third ply
		data{
			args: args{
				boardInFEN: "kn6/n6q/PP6/8" +
					"/8/8/7P/7K",
				color:       models.White,
				maximalDeep: 3,
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start:  models.Position{1, 5},
					Finish: models.Position{1, 6},
				},
				Score: -4,
			},
			wantErr: nil,
		},
	} {
		cache := caches.NewStringHashingCache(
			1e6,
			uci.EncodePieceStorage,
		)

		// increase the limit,
		// because the iterative searcher
		// discards a result
		// of the last iteration,
		// if it's not the only one
		if data.args.maximalDeep > 1 {
			data.args.maximalDeep++
		}

		gotMove, gotErr := iterativeSearch(
			cache,
			data.args.boardInFEN,
			data.args.color,
			data.args.maximalDeep,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Logf(
				"#%d:\ngot:  %v\nwant: %v",
				index,
				gotMove,
				data.wantMove,
			)

			test.Fail()
		}
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Logf(
				"#%d:\ngot:  %v\nwant: %v",
				index,
				gotErr,
				data.wantErr,
			)

			test.Fail()
		}
	}
}
