package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// MoveSearcherFactory ...
type MoveSearcherFactory func(
	terminator terminators.SearchTerminator,
) MoveSearcher

// ParallelSearcher ...
type ParallelSearcher struct {
	concurrency int
	factory     MoveSearcherFactory
}

// NewParallelSearcher ...
func NewParallelSearcher(
	concurrency int,
	factory MoveSearcherFactory,
) *ParallelSearcher {
	return &ParallelSearcher{
		concurrency: concurrency,
		factory:     factory,
	}
}

// SetSearcher ...
func (ParallelSearcher) SetSearcher(
	innerSearcher MoveSearcher,
) {
}

// SetTerminator ...
func (ParallelSearcher) SetTerminator(
	terminator terminators.SearchTerminator,
) {
}

// SearchMove ...
func (searcher ParallelSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	return moves.ScoredMove{}, nil
}
