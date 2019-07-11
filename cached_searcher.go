package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// Cache ...
type Cache interface {
	Get(
		storage models.PieceStorage,
		color models.Color,
	) (data moves.FailedMove, ok bool)
	Set(
		storage models.PieceStorage,
		color models.Color,
		data moves.FailedMove,
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
	innerSearcher MoveSearcher,
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
func (searcher CachedSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds Bounds,
) (moves.ScoredMove, error) {
	data, ok := searcher.cache.
		Get(storage, color)
	if ok {
		return data.Move, data.Error
	}

	move, err := searcher.searcher.
		SearchMove(storage, color, deep, bounds)
	data = moves.FailedMove{move, err}
	searcher.cache.Set(storage, color, data)

	return move, err
}
