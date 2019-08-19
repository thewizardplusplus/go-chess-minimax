package caches

import (
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
