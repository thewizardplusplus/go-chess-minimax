package chessminimax

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewParallelSearcher(
	test *testing.T,
) {
	factory := func(
		terminator terminators.SearchTerminator,
	) MoveSearcher {
		panic("not implemented")
	}
	searcher :=
		NewParallelSearcher(10, factory)

	if searcher.concurrency != 10 {
		test.Fail()
	}

	gotFactory := reflect.
		ValueOf(searcher.factory).
		Pointer()
	wantFactory := reflect.
		ValueOf(factory).
		Pointer()
	if gotFactory != wantFactory {
		test.Fail()
	}
}

func TestParallelSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		concurrency int
		factory     MoveSearcherFactory
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
				concurrency: 10,
				factory: func(
					terminator terminators.SearchTerminator,
				) MoveSearcher {
					_, ok := terminator.(*terminators.ParallelTerminator)
					if !ok {
						test.Fail()
					}

					return MockMoveSearcher{
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
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
								Score: 2.3,
							}
							return move,
								errors.New("dummy")
						},
					}
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
						File: 1,
						Rank: 2,
					},
					Finish: models.Position{
						File: 3,
						Rank: 4,
					},
				},
				Score: 2.3,
			},
			wantErr: true,
		},
	} {
		searcher := ParallelSearcher{
			concurrency: data.fields.concurrency,
			factory:     data.fields.factory,
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
