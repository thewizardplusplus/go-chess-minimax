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
	buffer := make(
		chan moves.FailedMove,
		searcher.concurrency,
	)
	terminator :=
		new(terminators.ParallelTerminator)
	for i := 0; i < searcher.concurrency; i++ {
		go func() {
			defer terminator.Terminate()

			searcher :=
				searcher.factory(terminator)
			move, err := searcher.SearchMove(
				storage,
				color,
				deep,
				bounds,
			)
			buffer <- moves.FailedMove{move, err}
		}()
	}

	move := <-buffer
	return move.Move, move.Error
}
