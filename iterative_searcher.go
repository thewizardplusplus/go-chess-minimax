package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// IterativeSearcher ...
type IterativeSearcher struct {
	MoveSearcher
}

// NewIterativeSearcher ...
func NewIterativeSearcher(
	innerSearcher MoveSearcher,
) *IterativeSearcher {
	searcher := &IterativeSearcher{
		MoveSearcher: innerSearcher,
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
	move, err := searcher.MoveSearcher.
		SearchMove(storage, color, deep, bounds)
	return move, err
}
