package generators

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// MoveGenerator ...
type MoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) []models.Move
}

// SafeMoveGenerator ...
type SafeMoveGenerator interface {
	MovesForColor(
		storage models.PieceStorage,
		color models.Color,
	) ([]models.Move, error)
}

// DefaultMoveGenerator ...
type DefaultMoveGenerator struct {
	innerGenerator MoveGenerator
}

// NewDefaultMoveGenerator ...
func NewDefaultMoveGenerator(
	innerGenerator MoveGenerator,
) DefaultMoveGenerator {
	return DefaultMoveGenerator{innerGenerator}
}

// MovesForColor ...
func (
	generator DefaultMoveGenerator,
) MovesForColor(
	storage models.PieceStorage,
	color models.Color,
) ([]models.Move, error) {
	moves := generator.innerGenerator.
		MovesForColor(storage, color)
	err := storage.CheckMoves(moves)
	if err != nil {
		return nil, err
	}

	return moves, nil
}
