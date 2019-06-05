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
	ErrCheck     = errors.New("check")
	ErrCheckmate = errors.New("checkmate")
)

var (
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
		MovesForColor(board, color)
	for _, move := range moves {
		piece, ok := board.Piece(move.Finish)
		if ok && piece.Kind() == models.King {
			return ScoredMove{}, ErrCheck
		}
	}

	bestMove := ScoredMove{
		Score: initialScore,
	}
	nextColor := color.Negative()
	for _, move := range moves {
		nextBoard := board.ApplyMove(move)
		scoredMove, err :=
			searcher.MoveSearcher.SearchMove(
				nextBoard,
				nextColor,
				deep+1,
			)
		switch err {
		case nil:
		case ErrCheck:
			continue
		default:
			return ScoredMove{}, err
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
