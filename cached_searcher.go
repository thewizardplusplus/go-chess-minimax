package chessminimax

import (
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// CachedSearcher ...
type CachedSearcher struct {
	MoveSearcher

	cache caches.Cache
}

// NewCachedSearcher ...
func NewCachedSearcher(
	cache caches.Cache,
	innerSearcher MoveSearcher,
) *CachedSearcher {
	searcher := &CachedSearcher{
		MoveSearcher: innerSearcher,

		cache: cache,
	}

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
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	data, ok := searcher.cache.
		Get(storage, color)
	if ok {
		return data.Move, data.Error
	}

	move, err := searcher.MoveSearcher.
		SearchMove(storage, color, deep, bounds)
	data = moves.FailedMove{move, err}
	searcher.cache.Set(storage, color, data)

	return move, err
}
