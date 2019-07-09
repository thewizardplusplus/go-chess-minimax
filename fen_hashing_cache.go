package chessminimax

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// FENHashingCache ...
type FENHashingCache map[string]CachedData

// Get ...
func (cache FENHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (data CachedData, ok bool) {
	fen, err := storage.ToFEN()
	if err != nil {
		return CachedData{}, false
	}

	data, ok = cache.data[fen]
	return applyColor(data, color), ok
}

// Set ...
func (cache FENHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	data CachedData,
) {
	fen, err := storage.ToFEN()
	if err != nil {
		return CachedData{}, false
	}

	cache.data[fen] = applyColor(data, color)
}

func applyColor(
	data CachedData,
	color models.Color,
) CachedData {
	if color == models.Black {
		data.Move.Score *= -1
	}

	return data
}
