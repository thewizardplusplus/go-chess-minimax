package chessminimax

import (
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockBoundedMoveSearcher struct {
	searchMove func(
		storage models.PieceStorage,
		color models.Color,
		deep int,
		alpha float64,
		beta float64,
	) (ScoredMove, error)
}

func (
	searcher MockBoundedMoveSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	alpha float64,
	beta float64,
) (ScoredMove, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(
		storage,
		color,
		deep,
		alpha,
		beta,
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
		alpha   float64
		beta    float64
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
						alpha float64,
						beta float64,
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
				color: models.White,
				deep:  2,
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
						alpha float64,
						beta float64,
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
				color: models.White,
				deep:  2,
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
						alpha float64,
						beta float64,
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
				color: models.White,
				deep:  2,
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
						alpha float64,
						beta float64,
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
				color: models.White,
				deep:  2,
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
						alpha float64,
						beta float64,
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
				color: models.White,
				deep:  2,
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
			data.args.alpha,
			data.args.beta,
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
