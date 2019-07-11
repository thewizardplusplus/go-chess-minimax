package caches

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPieceStorage struct {
	toFEN func() (string, error)
}

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

func (
	storage MockPieceStorage,
) ToFEN() (string, error) {
	if storage.toFEN == nil {
		panic("not implemented")
	}

	return storage.toFEN()
}

func TestFENHashingCacheGet(
	test *testing.T,
) {
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		cache     FENHashingCache
		args      args
		wantCache FENHashingCache
		wantData  moves.FailedMove
		wantOk    bool
	}

	for _, data := range []data{
		data{
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "", errors.New("dummy")
					},
				},
				color: models.White,
			},
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			wantData: moves.FailedMove{},
			wantOk:   false,
		},
		data{
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "fen #1", nil
					},
				},
				color: models.White,
			},
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			wantData: moves.FailedMove{
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
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.Black,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "fen #1", nil
					},
				},
				color: models.Black,
			},
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.Black,
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
			wantData: moves.FailedMove{
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
	} {
		gotData, gotOk := data.cache.Get(
			data.args.storage,
			data.args.color,
		)

		if !reflect.DeepEqual(
			data.cache,
			data.wantCache,
		) {
			test.Fail()
		}

		if !reflect.DeepEqual(
			gotData,
			data.wantData,
		) {
			test.Fail()
		}

		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}

func TestFENHashingCacheSet(
	test *testing.T,
) {
	type args struct {
		storage models.PieceStorage
		color   models.Color
		data    moves.FailedMove
	}
	type data struct {
		cache     FENHashingCache
		args      args
		wantCache FENHashingCache
	}

	for _, data := range []data{
		data{
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "", errors.New("dummy")
					},
				},
				color: models.White,
			},
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "fen #3", nil
					},
				},
				color: models.White,
				data: moves.FailedMove{
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
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #3",
					Color:      models.White,
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
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "fen #2", nil
					},
				},
				color: models.White,
				data: moves.FailedMove{
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
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #2",
					Color:      models.White,
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
			cache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.Black,
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
			args: args{
				storage: MockPieceStorage{
					toFEN: func() (string, error) {
						return "fen #1", nil
					},
				},
				color: models.Black,
				data: moves.FailedMove{
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
			wantCache: FENHashingCache{
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.White,
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
				FENHashKey{
					BoardInFEN: "fen #1",
					Color:      models.Black,
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
		data.cache.Set(
			data.args.storage,
			data.args.color,
			data.args.data,
		)

		if !reflect.DeepEqual(
			data.cache,
			data.wantCache,
		) {
			test.Fail()
		}
	}
}
