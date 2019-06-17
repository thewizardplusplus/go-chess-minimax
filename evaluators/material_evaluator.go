package evaluators

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

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

// MaterialEvaluator ...
type MaterialEvaluator struct{}

var (
	// are based on an evaluation function
	// of Claude Shannon
	pieceWeights = map[models.Kind]float64{
		models.King:   200,
		models.Queen:  9,
		models.Rook:   5,
		models.Bishop: 3,
		models.Knight: 3,
		models.Pawn:   1,
	}
)

// EvaluateBoard ...
func (
	evaluator MaterialEvaluator,
) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	var score float64
	for _, piece := range storage.Pieces() {
		weight := pieceWeights[piece.Kind()]
		colorSign := colorSign(piece, color)
		score += weight * colorSign
	}

	return score
}

func colorSign(
	piece models.Piece,
	color models.Color,
) float64 {
	if piece.Color() != color {
		return -1
	}

	return 1
}
