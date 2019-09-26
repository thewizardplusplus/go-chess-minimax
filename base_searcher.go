package chessminimax

import (
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
)

// BaseSearcher ...
type BaseSearcher struct {
	searcher   MoveSearcher
	terminator terminators.SearchTerminator
}

// SetSearcher ...
func (searcher *BaseSearcher) SetSearcher(
	innerSearcher MoveSearcher,
) {
	searcher.searcher = innerSearcher
}

// SetTerminator ...
func (searcher *BaseSearcher) SetTerminator(
	terminator terminators.SearchTerminator,
) {
	searcher.terminator = terminator
}
