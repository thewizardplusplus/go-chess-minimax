package caches

import (
	"container/list"
	"errors"
	"reflect"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
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
	maximalSize := int(1e6)
	cache := NewStringHashingCache(
		maximalSize,
		uci.EncodePieceStorage,
	)

	if cache.buckets == nil {
		test.Fail()
	}

	if cache.queue == nil {
		test.Fail()
	}

	if cache.maximalSize != maximalSize {
		test.Fail()
	}

	gotStringer := reflect.
		ValueOf(cache.stringer).
		Pointer()
	wantStringer := reflect.
		ValueOf(uci.EncodePieceStorage).
		Pointer()
	if gotStringer != wantStringer {
		test.Fail()
	}
}

func TestStringHashingCacheGet(
	test *testing.T,
) {
	type fields struct {
		buckets  bucketGroup
		queue    *list.List
		stringer Stringer
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
	}
	type data struct {
		fields    fields
		args      args
		wantQueue *list.List
		wantMove  moves.FailedMove
		wantOk    bool
	}

	for _, data := range []data{
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets: buckets,
					queue:   queue,
					stringer: func(
						storage models.PieceStorage,
					) string {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}

						return "key #2"
					},
				}
			}(),
			args: args{
				storage: MockPieceStorage{},
				color:   models.Black,
			},
			wantQueue: func() *list.List {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				queue := list.New()
				queue.PushBack(
					bucket{keyTwo, moveTwo},
				)
				queue.PushBack(
					bucket{keyOne, moveOne},
				)

				return queue
			}(),
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
					Score: 2.3,
				},
				Error: errors.New("dummy #2"),
			},
			wantOk: true,
		},
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets: buckets,
					queue:   queue,
					stringer: func(
						storage models.PieceStorage,
					) string {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}

						return "key #3"
					},
				}
			}(),
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
			},
			wantQueue: func() *list.List {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				queue := list.New()
				queue.PushBack(
					bucket{keyOne, moveOne},
				)
				queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return queue
			}(),
			wantMove: moves.FailedMove{},
			wantOk:   false,
		},
	} {
		cache := StringHashingCache{
			buckets:  data.fields.buckets,
			queue:    data.fields.queue,
			stringer: data.fields.stringer,
		}
		gotMove, gotOk := cache.Get(
			data.args.storage,
			data.args.color,
		)

		if !reflect.DeepEqual(
			data.fields.queue,
			data.wantQueue,
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
		buckets     bucketGroup
		queue       *list.List
		maximalSize int
		stringer    Stringer
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		move    moves.FailedMove
	}
	type data struct {
		fields     fields
		args       args
		wantFields fields
	}

	for _, data := range []data{
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets:     buckets,
					queue:       queue,
					maximalSize: 10,
					stringer: func(
						storage models.PieceStorage,
					) string {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}

						return "key #3"
					},
				}
			}(),
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantFields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}
				keyThree :=
					key{"key #3", models.White}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}
				moveThree := moves.FailedMove{
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyThree] = queue.PushBack(
					bucket{keyThree, moveThree},
				)
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets:     buckets,
					queue:       queue,
					maximalSize: 2,
					stringer: func(
						storage models.PieceStorage,
					) string {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}

						return "key #3"
					},
				}
			}(),
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantFields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyThree :=
					key{"key #3", models.White}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveThree := moves.FailedMove{
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyThree] = queue.PushBack(
					bucket{keyThree, moveThree},
				)
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 2.3,
					},
					Error: errors.New("dummy #2"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)

				return fields{
					buckets:     buckets,
					queue:       queue,
					maximalSize: 10,
					stringer: func(
						storage models.PieceStorage,
					) string {
						_, ok :=
							storage.(MockPieceStorage)
						if !ok {
							test.Fail()
						}

						return "key #2"
					},
				}
			}(),
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				},
			},
			wantFields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				moveOne := moves.FailedMove{
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
						Score: 1.2,
					},
					Error: errors.New("dummy #1"),
				}
				moveTwo := moves.FailedMove{
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
						Score: 4.2,
					},
					Error: errors.New("dummy #3"),
				}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyTwo] = queue.PushBack(
					bucket{keyTwo, moveTwo},
				)
				buckets[keyOne] = queue.PushBack(
					bucket{keyOne, moveOne},
				)

				return fields{
					buckets: buckets,
					queue:   queue,
				}
			}(),
		},
	} {
		cache := StringHashingCache{
			buckets:     data.fields.buckets,
			queue:       data.fields.queue,
			maximalSize: data.fields.maximalSize,
			stringer:    data.fields.stringer,
		}
		cache.Set(
			data.args.storage,
			data.args.color,
			data.args.move,
		)

		if !reflect.DeepEqual(
			data.fields.buckets,
			data.wantFields.buckets,
		) {
			test.Fail()
		}

		if !reflect.DeepEqual(
			data.fields.queue,
			data.wantFields.queue,
		) {
			test.Fail()
		}
	}
}

func TestStringHashingCacheMakeKey(
	test *testing.T,
) {
	cache := StringHashingCache{
		stringer: func(
			storage models.PieceStorage,
		) string {
			_, ok := storage.(MockPieceStorage)
			if !ok {
				test.Fail()
			}

			return "key"
		},
	}
	got := cache.makeKey(
		MockPieceStorage{},
		models.White,
	)

	want := key{"key", models.White}
	if !reflect.DeepEqual(got, want) {
		test.Fail()
	}
}

func TestStringHashingCacheGetElement(
	test *testing.T,
) {
	type fields struct {
		buckets     bucketGroup
		queue       *list.List
		wantElement *list.Element
	}
	type args struct {
		key key
	}
	type data struct {
		fields    fields
		args      args
		wantQueue *list.List
		wantOk    bool
	}

	for _, data := range []data{
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(1)
				buckets[keyTwo] = queue.PushBack(2)

				return fields{
					buckets:     buckets,
					queue:       queue,
					wantElement: buckets[keyTwo],
				}
			}(),
			args: args{
				key: key{"key #2", models.Black},
			},
			wantQueue: func() *list.List {
				queue := list.New()
				queue.PushBack(2)
				queue.PushBack(1)

				return queue
			}(),
			wantOk: true,
		},
		data{
			fields: func() fields {
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				buckets := make(bucketGroup)
				queue := list.New()
				buckets[keyOne] = queue.PushBack(1)
				buckets[keyTwo] = queue.PushBack(2)

				return fields{
					buckets:     buckets,
					queue:       queue,
					wantElement: nil,
				}
			}(),
			args: args{
				key: key{"key #3", models.White},
			},
			wantQueue: func() *list.List {
				queue := list.New()
				queue.PushBack(1)
				queue.PushBack(2)

				return queue
			}(),
			wantOk: false,
		},
	} {
		cache := StringHashingCache{
			buckets: data.fields.buckets,
			queue:   data.fields.queue,
		}
		gotElement, gotOk :=
			cache.getElement(data.args.key)

		if !reflect.DeepEqual(
			data.fields.queue,
			data.wantQueue,
		) {
			test.Fail()
		}

		if !reflect.DeepEqual(
			gotElement,
			data.fields.wantElement,
		) {
			test.Fail()
		}

		if gotOk != data.wantOk {
			test.Fail()
		}
	}
}
