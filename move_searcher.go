package chessminimax

import (
	"errors"
	"math"

	models "github.com/thewizardplusplus/go-chess-models"
)

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminate(
		board models.Board,
		deep int,
	) bool
}

// BoardEvaluator ...
type BoardEvaluator interface {
	EvaluateBoard(
		board models.Board,
		color models.Color,
	) float64
}

// MoveGenerator ...
type MoveGenerator interface {
	GenerateMoves(
		board models.Board,
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
		board models.Board,
		color models.Color,
		deep int,
	) (ScoredMove, error)
}

// DefaultMoveSearcher ...
type DefaultMoveSearcher struct {
	SearchTerminator SearchTerminator
	BoardEvaluator   BoardEvaluator
	MoveGenerator    MoveGenerator
	MoveSearcher     MoveSearcher
}

// ...
var (
	ErrNoKing    = errors.New("no king")
	ErrCheckmate = errors.New("checkmate")

	initialScore = math.Inf(-1)
)

// SearchMove ...
func (
	searcher DefaultMoveSearcher,
) SearchMove(
	board models.Board,
	color models.Color,
	deep int,
) (ScoredMove, error) {
	ok := searcher.SearchTerminator.
		IsSearchTerminate(board, deep)
	if ok {
		score := searcher.BoardEvaluator.
			EvaluateBoard(board, color)
		return ScoredMove{Score: score}, nil
	}

	moves := searcher.MoveGenerator.
		GenerateMoves(board, color)
	nextColor := color.Negative()
	for _, move := range moves {
		nextBoard := board.ApplyMove(move)
		if !hasKing(nextBoard, nextColor) {
			return ScoredMove{}, ErrNoKing
		}
	}

	bestMove := ScoredMove{
		Score: initialScore,
	}
	for _, move := range moves {
		nextBoard := board.ApplyMove(move)
		scoredMove, err :=
			searcher.MoveSearcher.SearchMove(
				nextBoard,
				nextColor,
				deep+1,
			)
		if err != nil {
			continue
		}

		score := -scoredMove.Score
		if bestMove.Score < score {
			bestMove = ScoredMove{move, score}
		}
	}
	if bestMove.Score == initialScore {
		return ScoredMove{}, ErrCheckmate
	}

	return bestMove, nil
}

func hasKing(
	board models.Board,
	color models.Color,
) bool {
	for _, piece := range board.Pieces() {
		if piece.Kind() == models.King &&
			piece.Color() == color {
			return true
		}
	}

	return false
}
