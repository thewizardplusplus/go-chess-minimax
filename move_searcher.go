package chessminimax

import (
	"errors"

	models "github.com/thewizardplusplus/go-chess-models"
)

// MoveGenerator ...
type MoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) []models.Move
}

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminate(deep int) bool
}

// BoardEvaluator ...
type BoardEvaluator interface {
	EvaluateBoard(
		storage models.PieceStorage,
		color models.Color,
	) float64
}

// MoveSearcher ...
type MoveSearcher interface {
	SearchMove(
		storage models.PieceStorage,
		color models.Color,
		deep int,
	) (ScoredMove, error)
}

// DefaultMoveSearcher ...
type DefaultMoveSearcher struct {
	generator  MoveGenerator
	terminator SearchTerminator
	evaluator  BoardEvaluator
	searcher   MoveSearcher
}

// ...
var (
	ErrCheck     = errors.New("check")
	ErrCheckmate = errors.New("checkmate")
	ErrDraw      = errors.New("draw")
)

// NewDefaultMoveSearcher ...
func NewDefaultMoveSearcher(
	generator MoveGenerator,
	terminator SearchTerminator,
	evaluator BoardEvaluator,
) *DefaultMoveSearcher {
	// instance must be created in a heap
	// so that it's possible to add
	// a reference to itself inside
	searcher := &DefaultMoveSearcher{
		generator:  generator,
		terminator: terminator,
		evaluator:  evaluator,
	}

	// use a reference to itself
	// for a recursion
	searcher.searcher = searcher

	return searcher
}

// SearchMove ...
func (
	searcher DefaultMoveSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
) (ScoredMove, error) {
	// check for a check should be first,
	// including before a termination check,
	// because a terminated evaluation
	// doesn't make sense for a check position
	moves := searcher.generator.
		MovesForColor(storage, color)
	err := storage.CheckMoves(moves)
	if err != nil {
		return ScoredMove{}, ErrCheck
	}

	ok := searcher.terminator.
		IsSearchTerminate(deep)
	if ok {
		score := searcher.evaluator.
			EvaluateBoard(storage, color)
		return ScoredMove{Score: score}, nil
	}

	bestMove := newScoredMove()
	var hasCheck bool
	for _, move := range moves {
		nextStorage := storage.ApplyMove(move)
		nextColor := color.Negative()
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				deep+1,
			)
		if err != nil {
			if err == ErrCheck {
				hasCheck = true
				continue
			}

			return ScoredMove{}, err
		}

		bestMove.update(scoredMove, move)
	}
	// no moves
	if !bestMove.isUpdated() {
		if hasCheck {
			return ScoredMove{}, ErrCheckmate
		}

		return ScoredMove{}, ErrDraw
	}

	return bestMove, nil
}
