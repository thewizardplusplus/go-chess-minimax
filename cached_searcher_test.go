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

func TestCachedSearcherSearchMove(
	test *testing.T,
) {
	type fields struct {
		searcher BoundedMoveSearcher
		cache    Cache
	}
	type args struct {
		storage models.PieceStorage
		color   models.Color
		deep    int
		bounds  Bounds
	}
	type data struct {
		fields   fields
		args     args
		wantMove ScoredMove
		wantErr  error
	}

	for _, data := range []data{} {
		searcher := CachedSearcher{
			cache: data.fields.cache,
		}
		searcher.searcher =
			data.fields.searcher

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
		if !reflect.DeepEqual(
			gotErr,
			data.wantErr,
		) {
			test.Fail()
		}
	}
}
