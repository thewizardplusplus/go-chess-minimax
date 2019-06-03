package chessminimax

import (
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
	) ScoredMove
}

// DefaultMoveSearcher ...
type DefaultMoveSearcher struct {
	SearchTerminator SearchTerminator
	BoardEvaluator   BoardEvaluator
	MoveGenerator    MoveGenerator
	MoveSearcher     MoveSearcher
}

// SearchMove ...
func (
	searcher DefaultMoveSearcher,
) SearchMove(
	board models.Board,
	color models.Color,
	deep int,
) ScoredMove {
	ok := searcher.SearchTerminator.
		IsSearchTerminate(board, deep)
	if ok {
		score := searcher.BoardEvaluator.
			EvaluateBoard(board, color)
		return ScoredMove{Score: score}
	}

	bestMove := ScoredMove{
		Score: math.Inf(+1),
	}
	moves := searcher.MoveGenerator.
		GenerateMoves(board, color)
	nextColor := negative(color)
	for _, move := range moves {
		nextBoard := board.ApplyMove(move)
		scoredMove :=
			searcher.MoveSearcher.SearchMove(
				nextBoard,
				nextColor,
				deep+1,
			)
		scoredMove.Score *= -1

		if bestMove.Score < scoredMove.Score {
			bestMove = scoredMove
		}
	}

	return bestMove
}

func negative(
	color models.Color,
) models.Color {
	if color == models.Black {
		return models.White
	}

	return models.Black
}
