package chessminimax

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// FENHashKey ...
type FENHashKey struct {
	BoardInFEN string
	Color      models.Color
}

// FENHashingCache ...
type FENHashingCache map[FENHashKey]moves.FailedMove

// Get ...
func (cache FENHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (data moves.FailedMove, ok bool) {
	fen, err := storage.ToFEN()
	if err != nil {
		return moves.FailedMove{}, false
	}

	data, ok = cache[FENHashKey{fen, color}]
	return data, ok
}

// Set ...
func (cache FENHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	data moves.FailedMove,
) {
	fen, err := storage.ToFEN()
	if err != nil {
		return
	}

	cache[FENHashKey{fen, color}] = data
}
