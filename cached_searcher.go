package chessminimax

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// CachedData ...
type CachedData struct {
	Move  ScoredMove
	Error error
}

// Cache ...
type Cache interface {
	Get(
		storage models.PieceStorage,
	) (data CachedData, ok bool)
	Set(
		storage models.PieceStorage,
		data CachedData,
	)
}

// CachedSearcher ...
type CachedSearcher struct {
	searcherHolder

	cache Cache
}

// NewCachedSearcher ...
func NewCachedSearcher(
	cache Cache,
	innerSearcher BoundedMoveSearcher,
) *CachedSearcher {
	searcher := &CachedSearcher{cache: cache}
	searcher.searcher = innerSearcher

	// set itself as an inner searcher
	// for passed one
	// in order to recursive calls
	// will be cached too
	innerSearcher.SetSearcher(searcher)

	return searcher
}

// SearchMove ...
func (
	searcher CachedSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds Bounds,
) (ScoredMove, error) {
	data, ok := searcher.cache.Get(storage)
	if ok {
		return data.Move, data.Error
	}

	move, err := searcher.searcher.
		SearchMove(storage, color, deep, bounds)
	if color == models.Black {
		move.Score *= -1
	}

	data = CachedData{move, err}
	searcher.cache.Set(storage, data)

	return move, err
}
