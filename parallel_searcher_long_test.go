// +build long

package chessminimax

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func TestParallelSearcher(test *testing.T) {
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
		// draw with checks
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
		// single profitable move
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
	} {
		shardedCache := caches.NewShardedCache(
			runtime.NumCPU(),
			func() caches.Cache {
				cache :=
					caches.NewStringHashingCache(
						1e6/runtime.NumCPU(),
						uci.EncodePieceStorage,
					)
				return caches.NewParallelCache(
					cache,
				)
			},
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

		gotMove, gotErr := parallelSearch(
			shardedCache,
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