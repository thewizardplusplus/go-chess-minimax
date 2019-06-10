package chessminimax

import (
	"errors"

	models "github.com/thewizardplusplus/go-chess-models"
)

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

// MoveGenerator ...
type MoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) []models.Move
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
	terminator SearchTerminator
	evaluator  BoardEvaluator
	generator  MoveGenerator
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
	terminator SearchTerminator,
	evaluator BoardEvaluator,
	generator MoveGenerator,
) *DefaultMoveSearcher {
	// instance must be created in a heap
	// so that it's possible to add
	// a reference to itself inside
	searcher := &DefaultMoveSearcher{
		terminator: terminator,
		evaluator:  evaluator,
		generator:  generator,
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
	ok := searcher.terminator.
		IsSearchTerminate(deep)
	if ok {
		score := searcher.evaluator.
			EvaluateBoard(storage, color)
		return ScoredMove{Score: score}, nil
	}

	moves := searcher.generator.
		MovesForColor(storage, color)
	err := storage.CheckMoves(moves)
	if err != nil {
		return ScoredMove{}, ErrCheck
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
