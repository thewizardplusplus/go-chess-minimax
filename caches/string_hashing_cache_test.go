package caches

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPieceStorage struct{}

func (
	storage MockPieceStorage,
) Size() models.Size {
	panic("not implemented")
}

func (storage MockPieceStorage) Piece(
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

func TestNewStringHashingCache(
	test *testing.T,
) {
	stringer := func(
		storage models.PieceStorage,
	) string {
		panic("not implemented")
	}
	cache := NewStringHashingCache(stringer)

	moves := make(moveGroup)
	if !reflect.DeepEqual(
		cache.moves,
		moves,
	) {
		test.Fail()
	}

	gotStringer := reflect.
		ValueOf(cache.stringer).
		Pointer()
	wantStringer := reflect.
		ValueOf(stringer).
		Pointer()
	if gotStringer != wantStringer {
		test.Fail()
	}
}

func TestStringHashingCacheGet(
	test *testing.T,
) {
	type fields struct {
		moves    moveGroup
		stringer Stringer
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		fields    fields
		args      args
		wantMoves moveGroup
		wantMove  moves.FailedMove
		wantOk    bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #2",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #1"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #2",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
			wantMove: moves.FailedMove{
				Move: moves.ScoredMove{
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
					Score: 2.3,
				},
				Error: errors.New("dummy #1"),
			},
			wantOk: true,
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #1",
						color:   models.Black,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #1"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.Black,
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #1",
					color:   models.Black,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
			wantMove: moves.FailedMove{
				Move: moves.ScoredMove{
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
				Error: errors.New("dummy #2"),
			},
			wantOk: true,
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #2",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #0"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #2",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
			wantMove: moves.FailedMove{},
			wantOk:   false,
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #2",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #1"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.Black,
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #2",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
			wantMove: moves.FailedMove{},
			wantOk:   false,
		},
	} {
		cache := StringHashingCache{
			moves:    data.fields.moves,
			stringer: data.fields.stringer,
		}
		gotMove, gotOk := cache.Get(
			data.args.storage,
			data.args.color,
		)

		if !reflect.DeepEqual(
			cache.moves,
			data.wantMoves,
		) {
			test.Fail()
		}

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}

		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}

func TestStringHashingCacheSet(
	test *testing.T,
) {
	type fields struct {
		moves    moveGroup
		stringer Stringer
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		move    moves.FailedMove
	}
	type data struct {
		fields    fields
		args      args
		wantMoves moveGroup
	}

	for _, data := range []data{
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #2",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #3"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				move: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #2",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
				key{
					storage: "fen #3",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #1"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.Black,
				move: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #1",
					color:   models.Black,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
					Error: errors.New("dummy #2"),
				},
			},
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #2",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #2"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				move: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #2",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
		},
		data{
			fields: fields{
				moves: moveGroup{
					key{
						storage: "fen #1",
						color:   models.White,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
							Score: 2.3,
						},
						Error: errors.New("dummy #1"),
					},
					key{
						storage: "fen #1",
						color:   models.Black,
					}: moves.FailedMove{
						Move: moves.ScoredMove{
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
						Error: errors.New("dummy #2"),
					},
				},
				stringer: func(
					storage models.PieceStorage,
				) string {
					_, ok :=
						storage.(MockPieceStorage)
					if !ok {
						test.Fail()
					}

					return "fen #1"
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.Black,
				move: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantMoves: moveGroup{
				key{
					storage: "fen #1",
					color:   models.White,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #1"),
				},
				key{
					storage: "fen #1",
					color:   models.Black,
				}: moves.FailedMove{
					Move: moves.ScoredMove{
						Move: models.Move{
							Start: models.Position{
								File: 9,
								Rank: 10,
							},
							Finish: models.Position{
								File: 11,
								Rank: 12,
							},
						},
						Score: 1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
		},
	} {
		cache := StringHashingCache{
			moves:    data.fields.moves,
			stringer: data.fields.stringer,
		}
		cache.Set(
			data.args.storage,
			data.args.color,
			data.args.move,
		)

		if !reflect.DeepEqual(
			cache.moves,
			data.wantMoves,
		) {
			test.Fail()
		}
	}
}
