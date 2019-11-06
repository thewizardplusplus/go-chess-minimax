package chessminimax

import (
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// OldCachedSearcher ...
type OldCachedSearcher struct {
	*SearcherSetter

	cache caches.Cache
}

// NewOldCachedSearcher ...
//
// It binds itself to the passed searcher.
func NewOldCachedSearcher(
	innerSearcher MoveSearcher,
	cache caches.Cache,
) OldCachedSearcher {
	searcher := OldCachedSearcher{
		SearcherSetter: new(SearcherSetter),

		cache: cache,
	}

	// set itself as an inner searcher
	// for passed one
	// in order to recursive calls
	// will be cached too
	innerSearcher.SetSearcher(searcher)
	searcher.SetSearcher(innerSearcher)

	return searcher
}

// SetTerminator ...
func (
	searcher OldCachedSearcher,
) SetTerminator(
	terminator terminators.SearchTerminator,
) {
	searcher.searcher.SetTerminator(terminator)
}

// SearchProgress ...
func (
	searcher OldCachedSearcher,
) SearchProgress(deep int) float64 {
	return searcher.searcher.
		SearchProgress(deep)
}

// SearchMove ...
func (searcher OldCachedSearcher) SearchMove(
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

	move, err := searcher.searcher.
		SearchMove(storage, color, deep, bounds)
	if !move.Move.IsZero() {
		data := moves.FailedMove{move, err}
		searcher.cache.Set(storage, color, data)
	}

	return move, err
}
