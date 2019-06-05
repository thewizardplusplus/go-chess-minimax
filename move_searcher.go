package chessminimax

import (
	"errors"
	"math"

	models "github.com/thewizardplusplus/go-chess-models"
)

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminate(
		storage models.PieceStorage,
		deep int,
	) bool
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

// ScoredMove ...
type ScoredMove struct {
	Move  models.Move
	Score float64
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

var (
	initialScore = math.Inf(-1)
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
		IsSearchTerminate(storage, deep)
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

	bestMove := ScoredMove{
		Score: initialScore,
	}
	var hasCheck bool
	nextColor := color.Negative()
	for _, move := range moves {
		nextStorage := storage.ApplyMove(move)
		scoredMove, err :=
			searcher.searcher.SearchMove(
				nextStorage,
				nextColor,
				deep+1,
			)
		switch err {
		case nil:
		case ErrCheck:
			hasCheck = true
			continue
		default:
			return ScoredMove{}, err
		}

		score := -scoredMove.Score
		if bestMove.Score < score {
			bestMove = ScoredMove{move, score}
		}
	}
	// no moves
	if bestMove.Score == initialScore {
		if hasCheck {
			return ScoredMove{}, ErrCheckmate
		}

		return ScoredMove{}, ErrDraw
	}

	return bestMove, nil
}
