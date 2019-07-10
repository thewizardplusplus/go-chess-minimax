package chessminimax

import (
	"errors"
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

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
		wantData  CachedData
		wantOk    bool
	}

	for _, data := range []data{
		data{
			cache: FENHashingCache{
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
			wantData: CachedData{},
			wantOk:   false,
		},
		data{
			cache: FENHashingCache{
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
			wantData: CachedData{
				Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
			wantData: CachedData{
				Move: ScoredMove{
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
					Score: -2.3,
				},
				Error: errors.New("dummy #1"),
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
		data    CachedData
	}
	type data struct {
		cache     FENHashingCache
		args      args
		wantCache FENHashingCache
	}

	for _, data := range []data{
		data{
			cache: FENHashingCache{
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				data: CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #3": CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				color: models.Black,
				data: CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				"fen #3": CachedData{
					Move: ScoredMove{
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
						Score: -1.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
		},
		data{
			cache: FENHashingCache{
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
				data: CachedData{
					Move: ScoredMove{
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
				"fen #1": CachedData{
					Move: ScoredMove{
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
				"fen #2": CachedData{
					Move: ScoredMove{
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
