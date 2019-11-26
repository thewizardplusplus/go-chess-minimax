package evaluators

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// MaterialEvaluator ...
type MaterialEvaluator struct{}

// EvaluateBoard ...
func (evaluator MaterialEvaluator) EvaluateBoard(
	storage models.PieceStorage,
	color models.Color,
) float64 {
	var score float64
	for _, piece := range storage.Pieces() {
		pieceWeight := pieceWeight(piece)
		colorSign := colorSign(piece, color)
		score += pieceWeight * colorSign
	}

	return score
}

// it's based on an evaluation function of Claude Shannon
func pieceWeight(piece models.Piece) float64 {
	var pieceWeight float64
	switch piece.Kind() {
	case models.King:
		pieceWeight = 200
	case models.Queen:
		pieceWeight = 9
	case models.Rook:
		pieceWeight = 5
	case models.Bishop:
		pieceWeight = 3
	case models.Knight:
		pieceWeight = 3
	case models.Pawn:
		pieceWeight = 1
	}

	return pieceWeight
}

func colorSign(piece models.Piece, color models.Color) float64 {
	if piece.Color() != color {
		return -1
	}

	return 1
}
