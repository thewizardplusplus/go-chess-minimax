package caches

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// Cache ...
type Cache interface {
	Get(
		storage models.PieceStorage,
		color models.Color,
	) (data moves.FailedMove, ok bool)
	Set(
		storage models.PieceStorage,
		color models.Color,
		data moves.FailedMove,
	)
}
