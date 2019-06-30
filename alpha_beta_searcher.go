package chessminimax

import (
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// BoundedMoveSearcher ...
type BoundedMoveSearcher interface {
	SearchMove(
		storage models.PieceStorage,
		color models.Color,
		deep int,
		alpha float64,
		beta float64,
	) (ScoredMove, error)
}

// AlphaBetaSearcher ...
type AlphaBetaSearcher struct {
	generator  MoveGenerator
	terminator terminators.SearchTerminator
	evaluator  evaluators.BoardEvaluator
	searcher   BoundedMoveSearcher
}

// NewAlphaBetaSearcher ...
func NewAlphaBetaSearcher(
	generator MoveGenerator,
	terminator terminators.SearchTerminator,
	evaluator evaluators.BoardEvaluator,
) *AlphaBetaSearcher {
	// instance must be created in a heap
	// so that it's possible to add
	// a reference to itself inside
	searcher := &AlphaBetaSearcher{
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
	searcher AlphaBetaSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	alpha float64,
	beta float64,
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
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				deep+1,
				-beta,
				-alpha,
			)
		if err == models.ErrKingCapture {
			hasCheck = true
			continue
		}

		score := -scoredMove.Score
		if score > alpha {
			alpha = score
		}
		if score >= beta {
			return ScoredMove{move, score}, nil
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
