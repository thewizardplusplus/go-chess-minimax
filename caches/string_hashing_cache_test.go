package caches

import (
	"container/list"
	"reflect"
	"testing"

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
				buckets := make(bucketGroup)
				queue := list.New()
				keyOne :=
					key{"key #1", models.White}
				buckets[keyOne] =
					queue.PushBack(keyOne)
				keyTwo :=
					key{"key #2", models.Black}
				buckets[keyTwo] =
					queue.PushBack(keyTwo)

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
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				queue := list.New()
				queue.PushBack(keyTwo)
				queue.PushBack(keyOne)

				return queue
			}(),
			wantOk: true,
		},
		data{
			fields: func() fields {
				buckets := make(bucketGroup)
				queue := list.New()
				keyOne :=
					key{"key #1", models.White}
				buckets[keyOne] =
					queue.PushBack(keyOne)
				keyTwo :=
					key{"key #2", models.Black}
				buckets[keyTwo] =
					queue.PushBack(keyTwo)

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
				keyOne :=
					key{"key #1", models.White}
				keyTwo :=
					key{"key #2", models.Black}

				queue := list.New()
				queue.PushBack(keyOne)
				queue.PushBack(keyTwo)

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
