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

// NewIterativeSearcher ...
func NewIterativeSearcher(
	terminator terminators.SearchTerminator,
	innerSearcher MoveSearcher,
) *IterativeSearcher {
	searcher := &IterativeSearcher{
		MoveSearcher: innerSearcher,

		terminator: terminator,
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
	var lastMove moves.FailedMove
	const startDeep = 1
	for deep := startDeep; ; deep++ {
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
				deep,
				bounds,
			)
		isTimeout := searcher.terminator.
			IsSearchTerminate(deep)
		if deep == startDeep || !isTimeout {
			lastMove = moves.FailedMove{move, err}
		}
		if isTimeout {
			break
		}
	}

	return lastMove.Move, lastMove.Error
}
