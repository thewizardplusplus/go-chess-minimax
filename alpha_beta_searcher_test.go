package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPieceStorage struct {
	appliedMove models.Move

	applyMove func(move models.Move) models.PieceStorage
}

func (storage MockPieceStorage) Size() models.Size {
	panic("not implemented")
}

func (storage MockPieceStorage) Piece(
	position models.Position,
) (piece models.Piece, ok bool) {
	panic("not implemented")
}

func (storage MockPieceStorage) Pieces() []models.Piece {
	panic("not implemented")
}

func (storage MockPieceStorage) ApplyMove(
	move models.Move,
) models.PieceStorage {
	if storage.applyMove == nil {
		panic("not implemented")
	}

	return storage.applyMove(move)
}

func (storage MockPieceStorage) CheckMove(move models.Move) error {
	panic("not implemented")
}

type MockMoveGenerator struct {
	movesForColor func(
		storage models.PieceStorage,
		color models.Color,
	) ([]models.Move, error)
}

func (generator MockMoveGenerator) MovesForColor(
	storage models.PieceStorage,
	color models.Color,
) ([]models.Move, error) {
	if generator.movesForColor == nil {
		panic("not implemented")
	}

	return generator.movesForColor(storage, color)
}

type MockBoardEvaluator struct {
	evaluateBoard func(storage models.PieceStorage, color models.Color) float64
}

func (evaluator MockBoardEvaluator) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	if evaluator.evaluateBoard == nil {
		panic("not implemented")
	}

	return evaluator.evaluateBoard(storage, color)
}

func TestNewAlphaBetaSearcher(test *testing.T) {
	var generator MockMoveGenerator
	var terminator MockSearchTerminator
	var evaluator MockBoardEvaluator
	searcher := NewAlphaBetaSearcher(generator, terminator, evaluator)

	if !reflect.DeepEqual(searcher.generator, generator) {
		test.Fail()
	}
	if !reflect.DeepEqual(searcher.terminator, terminator) {
		test.Fail()
	}
	if !reflect.DeepEqual(searcher.evaluator, evaluator) {
		test.Fail()
	}

	// check a reference to itself
	if !reflect.DeepEqual(searcher.searcher, searcher) {
		test.Fail()
	}
}

func TestAlphaBetaSearcherSearchMove(test *testing.T) {
	type fields struct {
		generator  MoveGenerator
		terminator terminators.SearchTerminator
		evaluator  evaluators.BoardEvaluator
		searcher   MoveSearcher
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
		wantErr  error
	}

	for _, data := range []data{
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						if _, ok := storage.(MockPieceStorage); !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return nil, models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{Alpha: -2e6, Beta: 3e6},
			},
			wantMove: moves.ScoredMove{},
			wantErr:  models.ErrKingCapture,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						if _, ok := storage.(MockPieceStorage); !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						moves := []models.Move{
							{
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
					isSearchTerminated: func(deep int) bool {
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
						if _, ok := storage.(MockPieceStorage); !ok {
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
				bounds:  moves.Bounds{Alpha: -2e6, Beta: 3e6},
			},
			wantMove: moves.ScoredMove{
				Score: 2.3,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}

						moves := []models.Move{
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6}) {
							test.Fail()
						}

						// all moves -> king capture
						return moves.ScoredMove{}, models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
			},
			wantMove: moves.ScoredMove{},
			wantErr:  ErrDraw,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}

						moves := []models.Move{
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
						// black color means a repeat call for checking, if a king
						// is under an attack
						if color == models.Black {
							err = models.ErrKingCapture
						}

						return moves, err
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6}) {
							test.Fail()
						}

						// all moves -> king capture
						return moves.ScoredMove{}, models.ErrKingCapture
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
			},
			wantMove: moves.ScoredMove{
				Score: evaluateCheckmate(2),
			},
			wantErr: ErrCheckmate,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6})
						checkTwo = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 4.2})
						if !checkOne && !checkTwo {
							test.Fail()
						}

						// move two -> king capture
						if checkTwo {
							return moves.ScoredMove{}, models.ErrKingCapture
						}

						// move one -> 4.2
						move := moves.ScoredMove{
							Score: 4.2,
						}
						return move, nil
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   -4.2,
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6})
						checkTwo = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 4.2})
						if !checkOne && !checkTwo {
							test.Fail()
						}

						var move moves.ScoredMove
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
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   -2.3,
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6}) {
							test.Fail()
						}

						var move moves.ScoredMove
						switch true {
						case checkOne:
							// move one -> less than the alpha
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
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   -2.3,
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6}) {
							test.Fail()
						}

						var move moves.ScoredMove
						switch true {
						case checkOne:
							// move one -> equal to the beta
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
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   3e6,
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}
						if !reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6}) {
							test.Fail()
						}

						var move moves.ScoredMove
						switch true {
						case checkOne:
							// move one -> greater than the beta
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
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   5e6,
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						mock, ok := storage.(MockPieceStorage)
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
							{
								Start: models.Position{
									File: 1,
									Rank: 2,
								},
								Finish: models.Position{
									File: 3,
									Rank: 4,
								},
							},
							{
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
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						checkOne := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						checkTwo := reflect.DeepEqual(storage, MockPieceStorage{
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
						})
						if !checkOne && !checkTwo {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						checkOne = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 2e6})
						checkTwo = reflect.DeepEqual(bounds, moves.Bounds{Alpha: -3e6, Beta: 4.2})
						if !checkOne && !checkTwo {
							test.Fail()
						}

						var move moves.ScoredMove
						var err error
						switch true {
						case checkOne:
							// move one -> 4.2
							move.Score = 4.2
						case checkTwo:
							// move two -> checkmate
							move.Score = evaluateCheckmate(3)
							err = ErrCheckmate
						}

						return move, err
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(move models.Move) models.PieceStorage {
						checkOne := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						})
						checkTwo := reflect.DeepEqual(move, models.Move{
							Start: models.Position{
								File: 5,
								Rank: 6,
							},
							Finish: models.Position{
								File: 7,
								Rank: 8,
							},
						})
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
				bounds: moves.Bounds{Alpha: -2e6, Beta: 3e6},
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
				Score:   -evaluateCheckmate(3),
				Quality: 0.25,
			},
			wantErr: nil,
		},
		{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) ([]models.Move, error) {
						if _, ok := storage.(MockPieceStorage); !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return nil, nil
					},
				},
				terminator: MockSearchTerminator{
					isSearchTerminated: func(deep int) bool {
						if deep != 2 {
							test.Fail()
						}

						return false
					},
					searchProgress: func(deep int) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.75
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{Alpha: -2e6, Beta: 3e6},
			},
			wantMove: moves.ScoredMove{},
			wantErr:  ErrDraw,
		},
	} {
		searcher := AlphaBetaSearcher{
			SearcherSetter: &SearcherSetter{
				searcher: data.fields.searcher,
			},
			TerminatorSetter: &TerminatorSetter{
				terminator: data.fields.terminator,
			},

			generator: data.fields.generator,
			evaluator: data.fields.evaluator,
		}

		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
			data.args.bounds,
		)

		if !reflect.DeepEqual(gotMove, data.wantMove) {
			test.Fail()
		}
		if !reflect.DeepEqual(gotErr, data.wantErr) {
			test.Fail()
		}
	}
}

func TestEvaluateCheckmate(test *testing.T) {
	scoreOne := evaluateCheckmate(2)
	scoreTwo := evaluateCheckmate(3)

	if scoreOne >= 0 {
		test.Fail()
	}
	if scoreTwo >= 0 {
		test.Fail()
	}
	if scoreTwo >= scoreOne {
		test.Fail()
	}
}
