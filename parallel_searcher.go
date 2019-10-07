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
	*TerminatorSetter

	concurrency int
	factory     MoveSearcherFactory
}

// NewParallelSearcher ...
func NewParallelSearcher(
	terminator terminators.SearchTerminator,
	concurrency int,
	factory MoveSearcherFactory,
) ParallelSearcher {
	searcher := ParallelSearcher{
		TerminatorSetter: new(TerminatorSetter),

		concurrency: concurrency,
		factory:     factory,
	}

	searcher.SetTerminator(terminator)

	return searcher
}

// SetSearcher ...
//
// It does nothing and is required
// only for correspondence
// to the MoveSearcher interface.
//
// It always panics.
func (ParallelSearcher) SetSearcher(
	innerSearcher MoveSearcher,
) {
	panic("not supported")
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
	terminator.Terminate()

	return move.Move, move.Error
}
