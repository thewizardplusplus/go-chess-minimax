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
//
// It binds itself to the passed searcher.
func NewCachedSearcher(
	innerSearcher MoveSearcher,
	cache caches.Cache,
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
	if !move.Move.IsZero() {
		data := moves.FailedMove{move, err}
		searcher.cache.Set(storage, color, data)
	}

	return move, err
}
