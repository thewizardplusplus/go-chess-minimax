package chessminimax

import (
	"errors"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// MoveGenerator ...
type MoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) ([]models.Move, error)
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
	generator  MoveGenerator
	terminator terminators.SearchTerminator
	evaluator  evaluators.BoardEvaluator
	searcher   MoveSearcher
}

// ...
var (
	ErrCheckmate = errors.New("checkmate")
	ErrDraw      = errors.New("draw")
)

// NewNegamaxSearcher ...
func NewNegamaxSearcher(
	generator MoveGenerator,
	terminator terminators.SearchTerminator,
	evaluator evaluators.BoardEvaluator,
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
		return ScoredMove{}, err
	}

	ok := searcher.terminator.
		IsSearchTerminate(deep)
	if ok {
		score := searcher.evaluator.
			EvaluateBoard(storage, color)
		return ScoredMove{Score: score}, nil
	}

	var hasCheck bool
	bestMove := newScoredMove()
	for _, move := range moves {
		nextStorage := storage.ApplyMove(move)
		nextColor := color.Negative()
		nextDeep := deep + 1
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				nextDeep,
			)
		if err == models.ErrKingCapture {
			hasCheck = true
			continue
		}

		bestMove.update(scoredMove, move)
	}
	// has a legal move
	if bestMove.isUpdated() {
		return bestMove, nil
	}

	// hasn't a legal move
	if hasCheck {
		// check, if a king is under an attack
		nextColor := color.Negative()
		_, err := searcher.generator.
			MovesForColor(storage, nextColor)
		if err != nil {
			score := evaluateCheckmate(deep)
			return ScoredMove{Score: score},
				ErrCheckmate
		}
	}

	// score of a draw is a null
	return ScoredMove{}, ErrDraw
}

// it evaluates a score of a checkmate
// for a current side, so its result
// should be negative
func evaluateCheckmate(deep int) float64 {
	score := 1e6 + float64(deep)
	return -score
}
