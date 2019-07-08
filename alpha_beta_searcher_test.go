package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockBoundedMoveSearcher struct {
	setInnerSearcher func(
		innerSearcher BoundedMoveSearcher,
	)
	searchMove func(
		storage models.PieceStorage,
		color models.Color,
		deep int,
		bounds Bounds,
	) (ScoredMove, error)
}

func (
	searcher MockBoundedMoveSearcher,
) SetInnerSearcher(
	innerSearcher BoundedMoveSearcher,
) {
	if searcher.setInnerSearcher == nil {
		panic("not implemented")
	}

	searcher.setInnerSearcher(innerSearcher)
}

func (
	searcher MockBoundedMoveSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds Bounds,
) (ScoredMove, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(
		storage,
		color,
		deep,
		bounds,
	)
}

func TestNewAlphaBetaSearcher(
	test *testing.T,
) {
	var generator MockMoveGenerator
	var terminator MockSearchTerminator
	var evaluator MockBoardEvaluator
	searcher := NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)

	if !reflect.DeepEqual(
		searcher.generator,
		generator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.terminator,
		terminator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.evaluator,
		evaluator,
	) {
		test.Fail()
	}

	// check a reference to itself
	if !reflect.DeepEqual(
		searcher.searcher,
		searcher,
	) {
		test.Fail()
	}
}

func TestAlphaBetaSearcherSetInnerSearcher(
	test *testing.T,
) {
	var generator MockMoveGenerator
	var terminator MockSearchTerminator
	var evaluator MockBoardEvaluator
	searcher := NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)

	var innerSearcher MockBoundedMoveSearcher
	searcher.SetInnerSearcher(innerSearcher)

	if !reflect.DeepEqual(
		searcher.generator,
		generator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.terminator,
		terminator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.evaluator,
		evaluator,
	) {
		test.Fail()
	}
	if !reflect.DeepEqual(
		searcher.searcher,
		innerSearcher,
	) {
		test.Fail()
	}
}

func TestAlphaBetaSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		generator  MoveGenerator
		terminator terminators.SearchTerminator
		evaluator  evaluators.BoardEvaluator
		searcher   BoundedMoveSearcher
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
		bounds  Bounds
	}
	type data struct {
		fields   fields
		args     args
		wantMove ScoredMove
		wantErr  error
	}

	for _, data := range []data{
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return nil,
							models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{},
			wantErr:  models.ErrKingCapture,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return true
					},
				},
				evaluator: MockBoardEvaluator{
					evaluateBoard: func(
						storage models.PieceStorage,
						color models.Color,
					) float64 {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return 2.3
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{Score: 2.3},
			wantErr:  nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						) {
							test.Fail()
						}

						// all moves -> king capture
						return ScoredMove{},
							models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}

						var err error
						// black color means
						// a repeat call for checking,
						// if a king is under an attack
						if color == models.Black {
							err = models.ErrKingCapture
						}

						return moves, err
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						) {
							test.Fail()
						}

						// all moves -> king capture
						return ScoredMove{},
							models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
				Score: evaluateCheckmate(2),
			},
			wantErr: ErrCheckmate,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						)
						checkTwo = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 4.2},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						// move two -> king capture
						if checkTwo {
							return ScoredMove{},
								models.ErrKingCapture
						}

						// move one -> 4.2
						move := ScoredMove{Score: 4.2}
						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: -4.2,
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						)
						checkTwo = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 4.2},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						var move ScoredMove
						switch true {
						case checkOne:
							// move one -> 4.2
							move.Score = 4.2
						case checkTwo:
							// move two -> 2.3
							move.Score = 2.3
						}

						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: -2.3,
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						) {
							test.Fail()
						}

						var move ScoredMove
						switch true {
						case checkOne:
							// move one ->
							//   less than the alpha
							move.Score = 5e6
						case checkTwo:
							// move two -> 2.3
							move.Score = 2.3
						}

						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: -2.3,
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						) {
							test.Fail()
						}

						var move ScoredMove
						switch true {
						case checkOne:
							// move one ->
							//   equal to the beta
							move.Score = -3e6
						case checkTwo:
							// move two -> 2.3
							move.Score = 2.3
						}

						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: 3e6,
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						) {
							test.Fail()
						}

						var move ScoredMove
						switch true {
						case checkOne:
							// move one ->
							//   greater than the beta
							move.Score = -5e6
						case checkTwo:
							// move two -> 2.3
							move.Score = 2.3
						}

						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: 5e6,
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						}
						return moves, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
				searcher: MockBoundedMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds Bounds,
					) (ScoredMove, error) {
						checkOne := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 1,
										Rank: 2,
									},
									Finish: models.Position{
										File: 3,
										Rank: 4,
									},
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							storage,
							MockPieceStorage{
								appliedMove: models.Move{
									Start: models.Position{
										File: 5,
										Rank: 6,
									},
									Finish: models.Position{
										File: 7,
										Rank: 8,
									},
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 2e6},
						)
						checkTwo = reflect.DeepEqual(
							bounds,
							Bounds{-3e6, 4.2},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						var move ScoredMove
						var err error
						switch true {
						case checkOne:
							// move one -> 4.2
							move.Score = 4.2
						case checkTwo:
							// move two -> checkmate
							move.Score =
								evaluateCheckmate(3)
							err = ErrCheckmate
						}

						return move, err
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						checkOne := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
						)
						checkTwo := reflect.DeepEqual(
							move,
							models.Move{
								Start: models.Position{
									File: 5,
									Rank: 6,
								},
								Finish: models.Position{
									File: 7,
									Rank: 8,
								},
							},
						)
						if !checkOne && !checkTwo {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
				},
				color:  models.White,
				deep:   2,
				bounds: Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{
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
				Score: -evaluateCheckmate(3),
			},
			wantErr: nil,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return nil, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminate: func(
						deep int,
					) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  Bounds{-2e6, 3e6},
			},
			wantMove: ScoredMove{},
			wantErr:  ErrDraw,
		},
	} {
		searcher := AlphaBetaSearcher{
			generator:  data.fields.generator,
			terminator: data.fields.terminator,
			evaluator:  data.fields.evaluator,
			searcher:   data.fields.searcher,
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
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Fail()
		}
	}
}
