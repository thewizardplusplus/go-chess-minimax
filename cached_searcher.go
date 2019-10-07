package chessminimax

import (
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// CachedSearcher ...
type CachedSearcher struct {
	*SearcherSetter

	cache caches.Cache
}

// NewCachedSearcher ...
//
// It binds itself to the passed searcher.
func NewCachedSearcher(
	innerSearcher MoveSearcher,
	cache caches.Cache,
) CachedSearcher {
	searcher := CachedSearcher{
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
func (searcher CachedSearcher) SetTerminator(
	terminator terminators.SearchTerminator,
) {
	searcher.searcher.SetTerminator(terminator)
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

	move, err := searcher.searcher.
		SearchMove(storage, color, deep, bounds)
	if !move.Move.IsZero() {
		data := moves.FailedMove{move, err}
		searcher.cache.Set(storage, color, data)
	}

	return move, err
}
