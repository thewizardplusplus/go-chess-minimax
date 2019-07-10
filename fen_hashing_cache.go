package chessminimax

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// FENHashKey ...
type FENHashKey struct {
	BoardInFEN string
	Color      models.Color
}

// FENHashingCache ...
type FENHashingCache map[FENHashKey]CachedData

// Get ...
func (cache FENHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (data CachedData, ok bool) {
	fen, err := storage.ToFEN()
	if err != nil {
		return CachedData{}, false
	}

	data, ok = cache[FENHashKey{fen, color}]
	return data, ok
}

// Set ...
func (cache FENHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	data CachedData,
) {
	fen, err := storage.ToFEN()
	if err != nil {
		return
	}

	cache[FENHashKey{fen, color}] = data
}
