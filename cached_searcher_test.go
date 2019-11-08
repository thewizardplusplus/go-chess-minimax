package chessminimax

import (
	"errors"
	"reflect"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
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

func TestNewCachedSearcher(
	test *testing.T,
) {
	innerSearcher := MockMoveSearcher{
		setSearcher: func(
			innerSearcher MoveSearcher,
		) {
			_, ok :=
				innerSearcher.(CachedSearcher)
			if !ok {
				test.Fail()
			}
		},
	}

	var cache MockCache
	searcher := NewCachedSearcher(
		innerSearcher,
		cache,
	)

	_, ok :=
		searcher.searcher.(MockMoveSearcher)
	if !ok {
		test.Fail()
	}

	if !reflect.DeepEqual(
		searcher.cache,
		cache,
	) {
		test.Fail()
	}
}

func TestCachedSearcherSetTerminator(
	test *testing.T,
) {
	var terminator MockSearchTerminator
	searcher := CachedSearcher{
		SearcherSetter: &SearcherSetter{
			searcher: MockMoveSearcher{
				setTerminator: func(
					innerTerminator terminators.SearchTerminator,
				) {
					if !reflect.DeepEqual(
						innerTerminator,
						terminator,
					) {
						test.Fail()
					}
				},
			},
		},
	}
	searcher.SetTerminator(terminator)
}

func TestCachedSearcherSearchProgress(
	test *testing.T,
) {
	searcher := CachedSearcher{
		SearcherSetter: &SearcherSetter{
			searcher: MockMoveSearcher{
				searchProgress: func(
					deep int,
				) float64 {
					if deep != 2 {
						test.Fail()
					}

					return 0.75
				},
			},
		},
	}
	got := searcher.SearchProgress(2)

	if got != 0.75 {
		test.Fail()
	}
}

func TestCachedSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		searcher MoveSearcher
		cache    caches.Cache
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
		wantErr  bool
	}

	for _, data := range []data{
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					searchProgress: func(
						deep int,
					) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.5
					},
				},
				cache: MockCache{
					get: func(
						storage models.PieceStorage,
						color models.Color,
					) (
						data moves.FailedMove,
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

						data = moves.FailedMove{
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
								Score:   2.3,
								Quality: 0.75,
							},
							Error: errors.New("dummy"),
						}
						return data, true
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
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
				Score:   2.3,
				Quality: 0.75,
			},
			wantErr: true,
		},
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					searchProgress: func(
						deep int,
					) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.5
					},
				},
				cache: MockCache{
					get: func(
						storage models.PieceStorage,
						color models.Color,
					) (
						data moves.FailedMove,
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

						data = moves.FailedMove{
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
								Score:   2.3,
								Quality: 0.5,
							},
							Error: errors.New("dummy"),
						}
						return data, true
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
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
				Score:   2.3,
				Quality: 0.5,
			},
			wantErr: true,
		},
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					searchProgress: func(
						deep int,
					) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.5
					},
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 2 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{-2e6, 3e6},
						) {
							test.Fail()
						}

						move := moves.ScoredMove{
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
						}
						return move, errors.New("dummy")
					},
				},
				cache: MockCache{
					get: func(
						storage models.PieceStorage,
						color models.Color,
					) (
						data moves.FailedMove,
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

						data = moves.FailedMove{
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
								Score:   2.3,
								Quality: 0.25,
							},
							Error: errors.New("dummy"),
						}
						return data, true
					},
					set: func(
						storage models.PieceStorage,
						color models.Color,
						data moves.FailedMove,
					) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						expectedData :=
							moves.FailedMove{
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
								Error: errors.New("dummy"),
							}
						if !reflect.DeepEqual(
							data,
							expectedData,
						) {
							test.Fail()
						}
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
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
				Score: 4.2,
			},
			wantErr: true,
		},
		data{
			fields: fields{
				searcher: MockMoveSearcher{
					searchProgress: func(
						deep int,
					) float64 {
						if deep != 2 {
							test.Fail()
						}

						return 0.5
					},
					searchMove: func(
						storage models.PieceStorage,
						color models.Color,
						deep int,
						bounds moves.Bounds,
					) (moves.ScoredMove, error) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}
						if deep != 2 {
							test.Fail()
						}
						if !reflect.DeepEqual(
							bounds,
							moves.Bounds{-2e6, 3e6},
						) {
							test.Fail()
						}

						move := moves.ScoredMove{
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
						}
						return move, errors.New("dummy")
					},
				},
				cache: MockCache{
					get: func(
						storage models.PieceStorage,
						color models.Color,
					) (
						data moves.FailedMove,
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

						return moves.FailedMove{}, false
					},
					set: func(
						storage models.PieceStorage,
						color models.Color,
						data moves.FailedMove,
					) {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}
						if color != models.White {
							test.Fail()
						}

						expectedData :=
							moves.FailedMove{
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
								Error: errors.New("dummy"),
							}
						if !reflect.DeepEqual(
							data,
							expectedData,
						) {
							test.Fail()
						}
					},
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
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
				Score: 4.2,
			},
			wantErr: true,
		},
	} {
		searcher := CachedSearcher{
			SearcherSetter: &SearcherSetter{
				searcher: data.fields.searcher,
			},

			cache: data.fields.cache,
		}

		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
			data.args.bounds,
		)

		if !reflect.DeepEqual(
			gotMove,
			data.wantMove,
		) {
			test.Fail()
		}

		hasErr := gotErr != nil
		if hasErr != data.wantErr {
			test.Fail()
		}
	}
}
