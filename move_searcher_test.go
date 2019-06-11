package chessminimax

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPieceStorage struct {
	appliedMove models.Move
	applyMove   func(
		move models.Move,
	) models.PieceStorage
	checkMoves func(
		moves []models.Move,
	) error
}

func (
	storage MockPieceStorage,
) Size() models.Size {
	panic("not implemented")
}

func (
	storage MockPieceStorage,
) Pieces() []models.Piece {
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

func (storage MockPieceStorage) CheckMove(
	move models.Move,
) error {
	panic("not implemented")
}

func (storage MockPieceStorage) CheckMoves(
	moves []models.Move,
) error {
	if storage.checkMoves == nil {
		panic("not implemented")
	}

	return storage.checkMoves(moves)
}

type MockMoveGenerator struct {
	movesForColor func(
		storage models.PieceStorage,
		color models.Color,
	) []models.Move
}

func (
	generator MockMoveGenerator,
) MovesForColor(
	storage models.PieceStorage,
	color models.Color,
) []models.Move {
	if generator.movesForColor == nil {
		panic("not implemented")
	}

	return generator.movesForColor(
		storage,
		color,
	)
}

type MockSearchTerminator struct {
	isSearchTerminate func(deep int) bool
}

func (
	terminator MockSearchTerminator,
) IsSearchTerminate(deep int) bool {
	if terminator.isSearchTerminate == nil {
		panic("not implemented")
	}

	return terminator.isSearchTerminate(deep)
}

type MockBoardEvaluator struct {
	evaluateBoard func(
		storage models.PieceStorage,
		color models.Color,
	) float64
}

func (
	evaluator MockBoardEvaluator,
) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	if evaluator.evaluateBoard == nil {
		panic("not implemented")
	}

	return evaluator.evaluateBoard(
		storage,
		color,
	)
}

type MockMoveSearcher struct {
	searchMove func(
		storage models.PieceStorage,
		color models.Color,
		deep int,
	) (ScoredMove, error)
}

func (searcher MockMoveSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
) (ScoredMove, error) {
	if searcher.searchMove == nil {
		panic("not implemented")
	}

	return searcher.searchMove(
		storage,
		color,
		deep,
	)
}

func TestMoveGeneratorInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		models.MoveGenerator{},
	)
	wantType := reflect.
		TypeOf((*MoveGenerator)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestMoveSearcherInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		DefaultMoveSearcher{},
	)
	wantType := reflect.
		TypeOf((*MoveSearcher)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestNewDefaultMoveSearcher(
	test *testing.T,
) {
	var generator MockMoveGenerator
	var terminator MockSearchTerminator
	var evaluator MockBoardEvaluator
	searcher := NewDefaultMoveSearcher(
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

func TestDefaultMoveSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		generator  MoveGenerator
		terminator SearchTerminator
		evaluator  BoardEvaluator
		searcher   MoveSearcher
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
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
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return models.ErrKingCapture
					},
				},
				color: models.White,
				deep:  2,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrCheck,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.checkMoves == nil {
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
				storage: MockPieceStorage{
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return nil
					},
				},
				color: models.White,
				deep:  2,
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
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
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

						// all moves -> check
						return ScoredMove{}, ErrCheck
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
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return nil
					},
				},
				color: models.White,
				deep:  2,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrCheckmate,
		},
		data{
			fields: fields{
				generator: MockMoveGenerator{
					movesForColor: func(
						storage models.PieceStorage,
						color models.Color,
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
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

						// move two -> check
						if checkTwo {
							return ScoredMove{}, ErrCheck
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
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return nil
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
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
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
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return nil
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
					) []models.Move {
						mock, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if mock.applyMove == nil {
							test.Fail()
						}
						if mock.checkMoves == nil {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						return []models.Move{
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
				searcher: MockMoveSearcher{
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
					) (ScoredMove, error) {
						expectedStorage :=
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
							}
						if !reflect.DeepEqual(
							storage,
							expectedStorage,
						) {
							test.Fail()
						}
						if color != models.Black {
							test.Fail()
						}
						if deep != 3 {
							test.Fail()
						}

						// move one -> checkmate
						var move ScoredMove
						return move, ErrCheckmate
					},
				},
			},
			args: args{
				storage: MockPieceStorage{
					applyMove: func(
						move models.Move,
					) models.PieceStorage {
						expectedMove := models.Move{
							Start: models.Position{
								File: 1,
								Rank: 2,
							},
							Finish: models.Position{
								File: 3,
								Rank: 4,
							},
						}
						if !reflect.DeepEqual(
							move,
							expectedMove,
						) {
							test.Fail()
						}

						return MockPieceStorage{
							appliedMove: move,
						}
					},
					checkMoves: func(
						moves []models.Move,
					) error {
						expectedMoves := []models.Move{
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
						if !reflect.DeepEqual(
							moves,
							expectedMoves,
						) {
							test.Fail()
						}

						return nil
					},
				},
				color: models.White,
				deep:  2,
			},
			wantMove: ScoredMove{},
			wantErr:  ErrCheckmate,
		},
	} {
		searcher := DefaultMoveSearcher{
			generator:  data.fields.generator,
			terminator: data.fields.terminator,
			evaluator:  data.fields.evaluator,
			searcher:   data.fields.searcher,
		}
		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
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
