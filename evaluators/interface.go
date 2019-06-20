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
