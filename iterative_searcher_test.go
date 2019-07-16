package chessminimax

import (
	"errors"
	"reflect"
	"testing"
	"time"

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

	maximalDuration := 5 * time.Second
	searcher := NewIterativeSearcher(
		innerSearcher,
		clock,
		maximalDuration,
	)

	_, ok := searcher.
		MoveSearcher.(MockMoveSearcher)
	if !ok {
		test.Fail()
	}

	gotClock := reflect.
		ValueOf(searcher.clock).
		Pointer()
	wantClock := reflect.
		ValueOf(clock).
		Pointer()
	if gotClock != wantClock {
		test.Fail()
	}

	if searcher.maximalDuration !=
		maximalDuration {
		test.Fail()
	}
}

func TestIterativeSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		searcher        MoveSearcher
		clock           terminators.Clock
		maximalDuration time.Duration
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
				clock:           clock,
				maximalDuration: 5 * time.Second,
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

			clock: data.fields.clock,
			maximalDuration: data.fields.
				maximalDuration,
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

func clock() time.Time {
	year, month, day := 2006, time.January, 2
	hour, minute, second := 15, 4, 5
	return time.Date(
		year, month, day,
		hour, minute, second,
		0,        // nanosecond
		time.UTC, // location
	)
}
