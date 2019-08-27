package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// IterativeSearcher ...
type IterativeSearcher struct {
	MoveSearcher

	terminator terminators.SearchTerminator
}

const (
	initialDeep = 1
)

// NewIterativeSearcher ...
func NewIterativeSearcher(
	innerSearcher MoveSearcher,
	terminator terminators.SearchTerminator,
) *IterativeSearcher {
	searcher := &IterativeSearcher{
		MoveSearcher: innerSearcher,

		terminator: terminator,
	}

	// set itself as an inner searcher
	// for passed one
	// in order to recursive calls
	// will be iterative too
	//innerSearcher.SetSearcher(searcher)

	return searcher
}

// SearchMove ...
func (searcher IterativeSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	var lastMove moves.FailedMove
	for deep := initialDeep; ; deep++ {
		searcher.MoveSearcher.SetTerminator(
			terminators.NewGroupTerminator(
				searcher.terminator,
				terminators.NewDeepTerminator(deep),
			),
		)

		move, err :=
			searcher.MoveSearcher.SearchMove(
				storage,
				color,
				0,
				bounds,
			)
		isTerminated := searcher.terminator.
			IsSearchTerminate(deep)
		if deep == initialDeep || !isTerminated {
			lastMove = moves.FailedMove{move, err}
		}
		// check at the loop end,
		// because there should be
		// at least one iteration
		if isTerminated {
			break
		}
	}

	return lastMove.Move, lastMove.Error
}
