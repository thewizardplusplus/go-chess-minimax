package chessminimax

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewParallelSearcher(
	test *testing.T,
) {
	var terminator MockSearchTerminator
	factory := func() MoveSearcher {
		panic("not implemented")
	}
	searcher := NewParallelSearcher(
		terminator,
		10,
		factory,
	)

	if searcher.concurrency != 10 {
		test.Fail()
	}

	gotFactory := reflect.
		ValueOf(searcher.factory).
		Pointer()
	wantFactory := reflect.
		ValueOf(factory).
		Pointer()
	if gotFactory != wantFactory {
		test.Fail()
	}
}

func TestParallelSearcherSetSearcher(
	test *testing.T,
) {
	var err interface{}
	func() {
		defer func() { err = recover() }()

		var innerSearcher MockMoveSearcher
		var searcher ParallelSearcher
		searcher.SetSearcher(innerSearcher)
	}()

	if err != "not supported" {
		test.Fail()
	}
}

func TestParallelSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		concurrency int
		factory     MoveSearcherFactory
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
		bounds  moves.Bounds
	}
	type data struct {
		fields fields
		args   args
	}

	var factoryWaiter *sync.WaitGroup
	var factoryCount uint64
	var searcherCount uint64
	var expectedMove moves.ScoredMove
	var expectedErr error
	var once sync.Once
	for _, data := range []data{
		data{
			fields: fields{
				concurrency: 10,
				factory: func() MoveSearcher {
					atomic.AddUint64(&factoryCount, 1)

					return MockMoveSearcher{
						searchMove: func(
							storage models.PieceStorage,
							color models.Color,
							deep int,
							bounds moves.Bounds,
						) (moves.ScoredMove, error) {
							defer factoryWaiter.Done()

							index := atomic.AddUint64(
								&searcherCount,
								1,
							)

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
										File: 1 + int(index),
										Rank: 2 + int(index),
									},
									Finish: models.Position{
										File: 3 + int(index),
										Rank: 4 + int(index),
									},
								},
								Score: 2.3,
							}
							err := errors.New("dummy")
							once.Do(func() {
								expectedMove = move
								expectedErr = err
							})

							return move, err
						},
					}
				},
			},
			args: args{
				storage: MockPieceStorage{},
				color:   models.White,
				deep:    2,
				bounds:  moves.Bounds{-2e6, 3e6},
			},
		},
	} {
		factoryWaiter = new(sync.WaitGroup)
		factoryWaiter.
			Add(data.fields.concurrency)

		atomic.StoreUint64(&factoryCount, 0)
		atomic.StoreUint64(&searcherCount, 0)

		searcher := ParallelSearcher{
			concurrency: data.fields.concurrency,
			factory:     data.fields.factory,
		}

		gotMove, gotErr := searcher.SearchMove(
			data.args.storage,
			data.args.color,
			data.args.deep,
			data.args.bounds,
		)
		factoryWaiter.Wait()

		gotFactoryCount :=
			atomic.LoadUint64(&factoryCount)
		gotSearcherCount :=
			atomic.LoadUint64(&searcherCount)

		if !reflect.DeepEqual(
			gotMove,
			expectedMove,
		) {
			test.Fail()
		}

		if gotErr != expectedErr {
			test.Fail()
		}

		if gotFactoryCount !=
			uint64(data.fields.concurrency) {
			test.Fail()
		}

		if gotSearcherCount !=
			uint64(data.fields.concurrency) {
			test.Fail()
		}
	}
}
