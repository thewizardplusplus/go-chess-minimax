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

/*func TestStringHashingCacheSet(
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
}*/
