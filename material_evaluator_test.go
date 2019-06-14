package chessminimax

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPiece struct {
	kind  func() models.Kind
	color func() models.Color
}

func (piece MockPiece) Kind() models.Kind {
	if piece.kind == nil {
		panic("not implemented")
	}

	return piece.kind()
}

func (
	piece MockPiece,
) Color() models.Color {
	if piece.color == nil {
		panic("not implemented")
	}

	return piece.color()
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
	board models.Board,
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
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.White
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.White
								},
							},
						}
					},
				},
				color: models.White,
			},
			want: 14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.White
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.Black
								},
							},
						}
					},
				},
				color: models.White,
			},
			want: 4,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.Black
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.Black
								},
							},
						}
					},
				},
				color: models.White,
			},
			want: -14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.White
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.White
								},
							},
						}
					},
				},
				color: models.Black,
			},
			want: -14,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.White
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.Black
								},
							},
						}
					},
				},
				color: models.Black,
			},
			want: -4,
		},
		data{
			args: args{
				storage: MockPieceStorage{
					pieces: func() []models.Piece {
						return []models.Piece{
							MockPiece{
								kind: func() models.Kind {
									return models.Queen
								},
								color: func() models.Color {
									return models.Black
								},
							},
							MockPiece{
								kind: func() models.Kind {
									return models.Rook
								},
								color: func() models.Color {
									return models.Black
								},
							},
						}
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
					color: func() models.Color {
						return models.Black
					},
				},
				color: models.White,
			},
			want: -1,
		},
		data{
			args: args{
				piece: MockPiece{
					color: func() models.Color {
						return models.White
					},
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
