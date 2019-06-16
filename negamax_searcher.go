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
//
// It should be a symmetric evaluation
// in relation to a side to move.
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

// NegamaxSearcher ...
type NegamaxSearcher struct {
	generator  SafeMoveGenerator
	terminator SearchTerminator
	evaluator  BoardEvaluator
	searcher   MoveSearcher
}

// ...
var (
	ErrCheckmate = errors.New("checkmate")
	ErrDraw      = errors.New("draw")
)

// NewNegamaxSearcher ...
func NewNegamaxSearcher(
	generator SafeMoveGenerator,
	terminator SearchTerminator,
	evaluator BoardEvaluator,
) *NegamaxSearcher {
	// instance must be created in a heap
	// so that it's possible to add
	// a reference to itself inside
	searcher := &NegamaxSearcher{
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
func (searcher NegamaxSearcher) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
) (ScoredMove, error) {
	// check for a check should be first,
	// including before a termination check,
	// because a terminated evaluation
	// doesn't make sense for a check position
	moves, err := searcher.generator.
		MovesForColor(storage, color)
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
	nextColor := color.Negative()
	var hasCheck bool
	for _, move := range moves {
		nextStorage := storage.ApplyMove(move)
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				deep+1,
			)
		if err == ErrCheck {
			hasCheck = true
			continue
		}

		bestMove.update(scoredMove, move)
	}
	// no moves
	if !bestMove.isUpdated() {
		if hasCheck {
			// check, if a king is under an attack
			_, err := searcher.generator.
				MovesForColor(storage, nextColor)
			if err != nil {
				score := evaluateCheckmate(deep)
				return ScoredMove{Score: score},
					ErrCheckmate
			}
		}

		return ScoredMove{}, ErrDraw
	}

	return bestMove, nil
}

// it evaluates a score of a checkmate
// for a current side, so its result
// should be negative
func evaluateCheckmate(deep int) float64 {
	score := 1e6 + float64(deep)
	return -score
}
