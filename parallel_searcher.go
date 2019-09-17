package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// ParallelSearcher ...
type ParallelSearcher struct {
	MoveSearcher

	terminator  terminators.SearchTerminator
	concurrency int
}

// NewParallelSearcher ...
func NewParallelSearcher(
	innerSearcher MoveSearcher,
	terminator terminators.SearchTerminator,
	concurrency int,
) *ParallelSearcher {
	return &ParallelSearcher{
		MoveSearcher: innerSearcher,

		terminator:  terminator,
		concurrency: concurrency,
	}
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
