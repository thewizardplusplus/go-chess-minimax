package chessminimax

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewIterativeSearcher(
	test *testing.T,
) {
	innerSearcher := MockMoveSearcher{
		setSearcher: func(
			innerSearcher MoveSearcher,
		) {
			mock, ok :=
				innerSearcher.(*IterativeSearcher)
			if !ok || mock == nil {
				test.Fail()
			}
		},
	}

	var terminator MockSearchTerminator
	searcher := NewIterativeSearcher(
		terminator,
		innerSearcher,
	)

	_, ok := searcher.
		MoveSearcher.(MockMoveSearcher)
	if !ok {
		test.Fail()
	}

	if !reflect.DeepEqual(
		searcher.terminator,
		terminator,
	) {
		test.Fail()
	}
}

func TestIterativeSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		searcher   MoveSearcher
		terminator terminators.SearchTerminator
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
		bounds  moves.Bounds
	}
	type data struct {
		fields   fields
		args     args
		wantMove moves.ScoredMove
		wantErr  bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 2 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{-2e6, 3e6},
						) {
							test.Fail()
						}

						move := moves.ScoredMove{
							Move: models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
							Score: 4.2,
						}
						return move, errors.New("dummy")
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						return true
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
			},
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start: models.Position{
						File: 5,
						Rank: 6,
					},
					Finish: models.Position{
						File: 7,
						Rank: 8,
					},
				},
				Score: 4.2,
			},
			wantErr: true,
		},
	} {
		searcher := IterativeSearcher{
			MoveSearcher: data.fields.searcher,

			terminator: data.fields.terminator,
		}

		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
			data.args.bounds,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}

		hasErr := gotErr != nil
		if hasErr != data.wantErr {
			test.Fail()
		}
	}
}
