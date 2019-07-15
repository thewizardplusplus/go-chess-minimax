package chessminimax

import (
  "github.com/thewizardplusplus/go-chess-minimax/terminators"
)

type baseSearcher struct {
  searcher   MoveSearcher
  terminator terminators.SearchTerminator
}

// SetSearcher ...
func (searcher *baseSearcher) SetSearcher(
  innerSearcher MoveSearcher,
) {
  searcher.searcher = innerSearcher
}

// SetTerminator ...
func (searcher *baseSearcher) SetTerminator(
  terminator terminators.SearchTerminator,
) {
  searcher.terminator = terminator
}
