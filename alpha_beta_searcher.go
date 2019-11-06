package chessminimax

import (
	"errors"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
)

// ...
var (
	ErrCheckmate = errors.New("checkmate")
	ErrDraw      = errors.New("draw")
)

// MoveGenerator ...
type MoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) ([]models.Move, error)
}

// AlphaBetaSearcher ...
type AlphaBetaSearcher struct {
	*SearcherSetter
	*TerminatorSetter

	generator MoveGenerator
	evaluator evaluators.BoardEvaluator
}

// NewAlphaBetaSearcher ...
func NewAlphaBetaSearcher(
	generator MoveGenerator,
	terminator terminators.SearchTerminator,
	evaluator evaluators.BoardEvaluator,
) AlphaBetaSearcher {
	// instance must be created in a heap
	// so that it's possible to add
	// a reference to itself inside
	searcher := AlphaBetaSearcher{
		SearcherSetter:   new(SearcherSetter),
		TerminatorSetter: new(TerminatorSetter),

		generator: generator,
		evaluator: evaluator,
	}

	// use a reference to itself
	// for a recursion
	searcher.SetSearcher(searcher)
	searcher.SetTerminator(terminator)

	return searcher
}

// SearchMove ...
func (
	searcher AlphaBetaSearcher,
) SearchMove(
	storage models.PieceStorage,
	color models.Color,
	deep int,
	bounds moves.Bounds,
) (moves.ScoredMove, error) {
	// check for a check should be first,
	// including before a termination check,
	// because a terminated evaluation
	// doesn't make sense for a check position
	moveGroup, err := searcher.generator.
		MovesForColor(storage, color)
	if err != nil {
		return moves.ScoredMove{}, err
	}

	ok := searcher.terminator.
		IsSearchTerminated(deep)
	if ok {
		score := searcher.evaluator.
			EvaluateBoard(storage, color)
		return moves.ScoredMove{Score: score},
			nil
	}

	var hasCheck bool
	bestMove := moves.NewScoredMove()
	moveQuality := 1 - searcher.terminator.
		SearchProgress(deep)
	for _, move := range moveGroup {
		nextStorage := storage.ApplyMove(move)
		nextColor := color.Negative()
		nextDeep := deep + 1
		nextBounds := bounds.Next()
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				nextDeep,
				nextBounds,
			)
		if err == models.ErrKingCapture {
			hasCheck = true
			continue
		}

		scoredMove, ok := bounds.Update(
			scoredMove,
			move,
			moveQuality,
		)
		if !ok {
			return scoredMove, nil
		}

		bestMove.Update(
			scoredMove,
			move,
			moveQuality,
		)
	}
	// has a legal move
	if bestMove.IsUpdated() {
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
			return moves.ScoredMove{Score: score},
				ErrCheckmate
		}
	}

	// score of a draw is a null
	return moves.ScoredMove{}, ErrDraw
}

// it evaluates a score of a checkmate
// for a current side, so its result
// should be negative
func evaluateCheckmate(deep int) float64 {
	score := 1e6 + float64(deep)
	return -score
}
