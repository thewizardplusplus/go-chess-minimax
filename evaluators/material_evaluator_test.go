package evaluators

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPieceStorage struct {
	pieces []models.Piece
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
	return storage.pieces
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

type MockPiece struct {
	kind  models.Kind
	color models.Color
}

func (piece MockPiece) Kind() models.Kind {
	return piece.kind
}

func (
	piece MockPiece,
) Color() models.Color {
	return piece.color
}

func (
	piece MockPiece,
) Position() models.Position {
	panic("not implemented")
}

func (piece MockPiece) ApplyPosition(
	position models.Position,
) models.Piece {
	panic("not implemented")
}

func (piece MockPiece) CheckMove(
	move models.Move,
	storage models.PieceStorage,
) bool {
	panic("not implemented")
}

func TestMaterialEvaluatorEvaluateBoard(
	test *testing.T,
) {
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.White,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.White,
						},
					},
				},
				color: models.White,
			},
			want: 14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.White,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.Black,
						},
					},
				},
				color: models.White,
			},
			want: 4,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.Black,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.Black,
						},
					},
				},
				color: models.White,
			},
			want: -14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.White,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.White,
						},
					},
				},
				color: models.Black,
			},
			want: -14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.White,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.Black,
						},
					},
				},
				color: models.Black,
			},
			want: -4,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: []models.Piece{
						MockPiece{
							kind:  models.Queen,
							color: models.Black,
						},
						MockPiece{
							kind:  models.Rook,
							color: models.Black,
						},
					},
				},
				color: models.Black,
			},
			want: 14,
		},
	} {
		var evaluator MaterialEvaluator
		got := evaluator.EvaluateBoard(
			data.args.storage,
			data.args.color,
		)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestPieceWeight(test *testing.T) {
	type args struct {
		piece models.Piece
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{
		data{
			args: args{
				piece: MockPiece{
					kind: models.Queen,
				},
			},
			want: 9,
		},
		data{
			args: args{
				piece: MockPiece{
					kind: models.Rook,
				},
			},
			want: 5,
		},
	} {
		got := pieceWeight(data.args.piece)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestColorSign(test *testing.T) {
	type args struct {
		piece models.Piece
		color models.Color
	}
	type data struct {
		args args
		want float64
	}

	for _, data := range []data{
		data{
			args: args{
				piece: MockPiece{
					color: models.Black,
				},
				color: models.White,
			},
			want: -1,
		},
		data{
			args: args{
				piece: MockPiece{
					color: models.White,
				},
				color: models.White,
			},
			want: 1,
		},
	} {
		got := colorSign(
			data.args.piece,
			data.args.color,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
