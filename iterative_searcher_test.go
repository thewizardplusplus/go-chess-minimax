package chessminimax

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewIterativeSearcher(test *testing.T) {
	var innerSearcher MockMoveSearcher
	var terminator MockSearchTerminator
	searcher := NewIterativeSearcher(innerSearcher, terminator)

	if !reflect.DeepEqual(searcher.searcher, innerSearcher) {
		test.Fail()
	}
	if !reflect.DeepEqual(searcher.terminator, terminator) {
		test.Fail()
	}
}

func TestIterativeSearcherSearchMove(test *testing.T) {
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
		{
			fields: fields{
				searcher: MockMoveSearcher{
					setTerminator: func(terminator terminators.SearchTerminator) {
						if _, ok := terminator.(terminators.GroupTerminator); !ok {
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

						if _, ok := storage.(MockPieceStorage); !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 0 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -2e6, Beta: 3e6}) {
							test.Fail()
						}

						move := moves.ScoredMove{
							Move: models.Move{
								Start: models.Position{
									File: testDeep + 1,
									Rank: testDeep + 2,
								},
								Finish: models.Position{
									File: testDeep + 3,
									Rank: testDeep + 4,
								},
							},
							Score: float64(testDeep + 5),
						}
						return move, errors.New("dummy")
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminated: func(deep int) bool {
						return deep == 1
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
		{
			fields: fields{
				searcher: MockMoveSearcher{
					setTerminator: func(terminator terminators.SearchTerminator) {
						if _, ok := terminator.(terminators.GroupTerminator); !ok {
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

						if _, ok := storage.(MockPieceStorage); !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 0 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -2e6, Beta: 3e6}) {
							test.Fail()
						}

						move := moves.ScoredMove{
							Move: models.Move{
								Start: models.Position{
									File: testDeep + 1,
									Rank: testDeep + 2,
								},
								Finish: models.Position{
									File: testDeep + 3,
									Rank: testDeep + 4,
								},
							},
							Score: float64(testDeep + 5),
						}
						return move, errors.New("dummy")
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminated: func(deep int) bool {
						return deep == 5
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
			SearcherSetter: &SearcherSetter{
				searcher: data.fields.searcher,
			},
			TerminatorSetter: &TerminatorSetter{
				terminator: data.fields.terminator,
			},
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
		if !reflect.DeepEqual(gotMove, data.wantMove) {
			test.Fail()
		}
		if hasErr := gotErr != nil; hasErr != data.wantErr {
			test.Fail()
		}
	}
}
