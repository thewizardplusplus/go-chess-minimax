package chessminimax

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// FENHashingCache ...
type FENHashingCache struct {
	data map[string]CachedData
}

// NewFENHashingCache ...
func NewFENHashingCache() FENHashingCache {
	data := make(map[string]CachedData)
	return FENHashingCache{data}
}

// Get ...
func (cache FENHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (data CachedData, ok bool) {
	data, ok = cache.data[storage.ToFEN()]
	return data, ok
}

// Set ...
func (cache FENHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	data CachedData,
) {
	cache.data[storage.ToFEN()] = data
}
