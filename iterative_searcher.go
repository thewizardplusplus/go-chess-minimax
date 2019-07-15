package chessminimax

import (
	"time"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// IterativeSearcher ...
type IterativeSearcher struct {
	MoveSearcher

	clock           terminators.Clock
	maximalDuration time.Duration
}

// NewIterativeSearcher ...
func NewIterativeSearcher(
	innerSearcher MoveSearcher,
	clock terminators.Clock,
	maximalDuration time.Duration,
) *IterativeSearcher {
	searcher := &IterativeSearcher{
		MoveSearcher: innerSearcher,

		clock:           clock,
		maximalDuration: maximalDuration,
	}

	// set itself as an inner searcher
	// for passed one
	// in order to recursive calls
	// will be iterative too
	innerSearcher.SetSearcher(searcher)

	return searcher
}

// SearchMove ...
func (searcher IterativeSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	timer := terminators.NewTimeTerminator(
		searcher.clock,
		searcher.maximalDuration,
	)
	isTimeout := func() bool {
		return timer.IsSearchTerminate(0)
	}

	var lastMove moves.FailedMove
	for deep := 1; isTimeout(); deep++ {
		searcher.MoveSearcher.SetTerminator(
			terminators.NewGroupTerminator(
				timer,
				terminators.NewDeepTerminator(deep),
			),
		)

		move, err :=
			searcher.MoveSearcher.SearchMove(
				storage,
				color,
				deep,
				bounds,
			)
		isFirstIteration :=
			!lastMove.Move.IsUpdated()
		if isFirstIteration || !isTimeout() {
			lastMove = moves.FailedMove{move, err}
		}
	}

	return lastMove.Move, lastMove.Error
}
