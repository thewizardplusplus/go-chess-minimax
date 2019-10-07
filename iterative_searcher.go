package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// IterativeSearcher ...
type IterativeSearcher struct {
	*SearcherSetter
	*TerminatorSetter
}

const (
	initialDeep = 1
)

// NewIterativeSearcher ...
func NewIterativeSearcher(
	innerSearcher MoveSearcher,
	terminator terminators.SearchTerminator,
) IterativeSearcher {
	searcher := IterativeSearcher{
		SearcherSetter:   new(SearcherSetter),
		TerminatorSetter: new(TerminatorSetter),
	}

	searcher.SetSearcher(innerSearcher)
	searcher.SetTerminator(terminator)

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
		searcher.searcher.SetTerminator(
			terminators.NewGroupTerminator(
				searcher.terminator,
				terminators.NewDeepTerminator(deep),
			),
		)

		move, err :=
			searcher.searcher.SearchMove(
				storage,
				color,
				0,
				bounds,
			)
		isTerminated := searcher.terminator.
			IsSearchTerminated(deep)
		if deep == initialDeep || !isTerminated {
			lastMove = moves.FailedMove{
				Move:  move,
				Error: err,
			}
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
