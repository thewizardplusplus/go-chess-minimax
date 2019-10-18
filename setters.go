package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// MoveSearcher ...
type MoveSearcher interface {
	SetSearcher(searcher MoveSearcher)
	SetTerminator(
		terminator terminators.SearchTerminator,
	)

	// It should return only following errors:
	// * models.ErrKingCapture;
	// * games.ErrCheckmate;
	// * games.ErrDraw.
	SearchMove(
		storage models.PieceStorage,
		color models.Color,
		deep int,
		bounds moves.Bounds,
	) (moves.ScoredMove, error)
}

// SearcherSetter ...
type SearcherSetter struct {
	searcher MoveSearcher
}

// SetSearcher ...
func (setter *SearcherSetter) SetSearcher(
	searcher MoveSearcher,
) {
	setter.searcher = searcher
}

// TerminatorSetter ...
type TerminatorSetter struct {
	terminator terminators.SearchTerminator
}

// SetTerminator ...
func (
	setter *TerminatorSetter,
) SetTerminator(
	terminator terminators.SearchTerminator,
) {
	setter.terminator = terminator
}
