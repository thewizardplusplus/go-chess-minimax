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
		innerSearcher,
		terminator,
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
		wantDeep int
		wantMove moves.ScoredMove
		wantErr  bool
	}

	var testDeep int
	for _, data := range []data{
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					setTerminator: func(terminator terminators.SearchTerminator) {
						_, ok := terminator.(terminators.GroupTerminator)
						if !ok {
							test.Fail()
						}
					},
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						defer func() { testDeep++ }()

						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != testDeep {
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
									File: deep + 1,
									Rank: deep + 2,
								},
								Finish: models.Position{
									File: deep + 3,
									Rank: deep + 4,
								},
							},
							Score: float64(deep + 5),
						}
						return move, errors.New("dummy")
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						return deep == 1
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
			},
			wantDeep: 2,
			wantMove: moves.ScoredMove{
				Move: models.Move{
					Start: models.Position{
						File: 2,
						Rank: 3,
					},
					Finish: models.Position{
						File: 4,
						Rank: 5,
					},
				},
				Score: float64(6),
			},
			wantErr: true,
		},
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					setTerminator: func(terminator terminators.SearchTerminator) {
						_, ok := terminator.(terminators.GroupTerminator)
						if !ok {
							test.Fail()
						}
					},
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						defer func() { testDeep++ }()

						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != testDeep {
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
									File: deep + 1,
									Rank: deep + 2,
								},
								Finish: models.Position{
									File: deep + 3,
									Rank: deep + 4,
								},
							},
							Score: float64(deep + 5),
						}
						return move, errors.New("dummy")
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						return deep == 5
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
			},
			wantDeep: 6,
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
				Score: float64(9),
			},
			wantErr: true,
		},
	} {
		testDeep = 1

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

		if testDeep != data.wantDeep {
			test.Fail()
		}

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
