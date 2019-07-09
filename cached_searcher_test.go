package chessminimax

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockCache struct {
	get func(
		storage models.PieceStorage,
		color models.Color,
	) (data CachedData, ok bool)
	set func(
		storage models.PieceStorage,
		color models.Color,
		data CachedData,
	)
}

func (cache MockCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (data CachedData, ok bool) {
	if cache.get == nil {
		panic("not implemented")
	}

	return cache.get(storage, color)
}

func (cache MockCache) Set(
	storage models.PieceStorage,
	color models.Color,
	data CachedData,
) {
	if cache.set == nil {
		panic("not implemented")
	}

	cache.set(storage, color, data)
}

func TestNewCachedSearcher(test *testing.T) {
	innerSearcher := MockBoundedMoveSearcher{
		setSearcher: func(
			innerSearcher BoundedMoveSearcher,
		) {
			mock, ok :=
				innerSearcher.(*CachedSearcher)
			if !ok || mock == nil {
				test.Fail()
			}
		},
	}

	var cache Cache
	searcher := NewCachedSearcher(
		cache,
		innerSearcher,
	)

	if !reflect.DeepEqual(
		searcher.cache,
		cache,
	) {
		test.Fail()
	}

	_, ok := searcher.
		searcher.(MockBoundedMoveSearcher)
	if !ok {
		test.Fail()
	}
}
