package caches

import (
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

type MockCache struct {
	get func(
		storage models.PieceStorage,
		color models.Color,
	) (move moves.FailedMove, ok bool)
	set func(
		storage models.PieceStorage,
		color models.Color,
		move moves.FailedMove,
	)
}

func (cache MockCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	if cache.get == nil {
		panic("not implemented")
	}

	return cache.get(storage, color)
}

func (cache MockCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	if cache.set == nil {
		panic("not implemented")
	}

	cache.set(storage, color, move)
}

func TestNewParallelCache(test *testing.T) {
	var innerCache MockCache
	cache := NewParallelCache(innerCache)

	if !reflect.DeepEqual(
		cache.innerCache,
		innerCache,
	) {
		test.Fail()
	}
}

func TestParallelCacheGet(test *testing.T) {
	type fields struct {
		innerCache Cache
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		fields   fields
		args     args
		wantMove moves.FailedMove
		wantOk   bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				innerCache: MockCache{
					get: func(
						storage models.PieceStorage,
						color models.Color,
					) (
						move moves.FailedMove,
						ok bool,
					) {
						_, ok =
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						move = moves.FailedMove{
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
							Error: errors.New("dummy"),
						}
						return move, true
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
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
				Error: errors.New("dummy"),
			},
			wantOk: true,
		},
	} {
		cache := ParallelCache{
			innerCache: data.fields.innerCache,
		}
		gotMove, gotOk := cache.Get(
			data.args.storage,
			data.args.color,
		)

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

func TestParallelCacheSet(test *testing.T) {
	type fields struct {
		innerCache Cache
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		move    moves.FailedMove
	}
	type data struct {
		fields fields
		args   args
	}

	for _, data := range []data{
		data{
			fields: fields{
				innerCache: MockCache{
					set: func(
						storage models.PieceStorage,
						color models.Color,
						move moves.FailedMove,
					) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						expectedMove :=
							moves.FailedMove{
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
								Error: errors.New("dummy"),
							}
						if !reflect.DeepEqual(
							move,
							expectedMove,
						) {
							test.Fail()
						}
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				move: moves.FailedMove{
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
					Error: errors.New("dummy"),
				},
			},
		},
	} {
		cache := ParallelCache{
			innerCache: data.fields.innerCache,
		}
		cache.Set(
			data.args.storage,
			data.args.color,
			data.args.move,
		)
	}
}
