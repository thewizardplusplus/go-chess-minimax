package caches

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// Stringer ...
type Stringer func(
	storage models.PieceStorage,
) string

type key struct {
	storage string
	color   models.Color
}

type moveGroup map[key]moves.FailedMove

// StringHashingCache ...
type StringHashingCache struct {
	moves    moveGroup
	stringer Stringer
}

// NewStringHashingCache ...
func NewStringHashingCache(
	stringer Stringer,
) StringHashingCache {
	moves := make(moveGroup)
	return StringHashingCache{moves, stringer}
}

// Get ...
func (cache StringHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	key := cache.makeKey(storage, color)
	move, ok = cache.moves[key]
	return move, ok
}

// Set ...
func (cache StringHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	key := cache.makeKey(storage, color)
	cache.moves[key] = move
}

func (cache StringHashingCache) makeKey(
	storage models.PieceStorage,
	color models.Color,
) key {
	text := cache.stringer(storage)
	return key{text, color}
}
