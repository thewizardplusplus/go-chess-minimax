package generators

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

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

type MockPieceStorage struct {
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
) Piece(
	position models.Position,
) (piece models.Piece, ok bool) {
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
	panic("not implemented")
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

func TestSafeMoveGeneratorInterface(
	test *testing.T,
) {
	gotType := reflect.TypeOf(
		DefaultMoveGenerator{},
	)
	wantType := reflect.
		TypeOf((*SafeMoveGenerator)(nil)).
		Elem()
	if !gotType.Implements(wantType) {
		test.Fail()
	}
}

func TestDefaultMoveGeneratorMovesForColor(
	test *testing.T,
) {
	type fields struct {
		innerGenerator MoveGenerator
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		fields    fields
		args      args
		wantMoves []models.Move
		wantErr   error
	}

	for _, data := range []data{
		data{
			fields: fields{
				innerGenerator: MockMoveGenerator{
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
			},
			wantMoves: nil,
			wantErr:   models.ErrKingCapture,
		},
		data{
			fields: fields{
				innerGenerator: MockMoveGenerator{
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

						return nil
					},
				},
				color: models.White,
			},
			wantMoves: []models.Move{
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
			},
			wantErr: nil,
		},
	} {
		innerGenerator :=
			data.fields.innerGenerator
		generator := DefaultMoveGenerator{
			innerGenerator: innerGenerator,
		}
		gotMoves, gotErr :=
			generator.MovesForColor(
				data.args.storage,
				data.args.color,
			)

		if !reflect.DeepEqual(
			gotMoves,
			data.wantMoves,
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
