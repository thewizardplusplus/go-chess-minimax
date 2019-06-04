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

type appliedMove struct {
	move  models.Move
	board models.Board
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

	var appliedMoves []appliedMove
	moves := searcher.MoveGenerator.
		GenerateMoves(board, color)
	for _, move := range moves {
		nextBoard := board.ApplyMove(move)
		if !hasKing(nextBoard) {
			return ScoredMove{}, ErrNoKing
		}

		appliedMoves = append(
			appliedMoves,
			appliedMove{move, nextBoard},
		)
	}

	bestMove := ScoredMove{
		Score: initialScore,
	}
	nextColor := negative(color)
	for _, move := range appliedMoves {
		scoredMove, err :=
			searcher.MoveSearcher.SearchMove(
				move.board,
				nextColor,
				deep+1,
			)
		if err != nil {
			continue
		}

		score := -scoredMove.Score
		if bestMove.Score < score {
			bestMove = ScoredMove{
				Move:  move.move,
				Score: score,
			}
		}
	}
	if bestMove.Score == initialScore {
		return ScoredMove{}, ErrCheckmate
	}

	return bestMove, nil
}

func hasKing(board models.Board) bool {
	for _, piece := range board.Pieces() {
		if piece.Kind() == models.King {
			return true
		}
	}

	return false
}

func negative(
	color models.Color,
) models.Color {
	if color == models.Black {
		return models.White
	}

	return models.Black
}
