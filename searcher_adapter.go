package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// SearcherAdapter ...
type SearcherAdapter struct {
	MoveSearcher
}

// SearchMove ...
func (adapter SearcherAdapter) SearchMove(
	storage models.PieceStorage,
	color models.Color,
) (models.Move, error) {
	move, err :=
		adapter.MoveSearcher.SearchMove(
			storage,
			color,
			0, // initial deep
			moves.NewBounds(),
		)
	return move.Move, err
}
