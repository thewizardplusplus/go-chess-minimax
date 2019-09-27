package caches

import (
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// CacheFactory ...
type CacheFactory func() Cache

// ShardedCache ...
type ShardedCache struct {
	shards []Cache
}

// NewShardedCache ...
func NewShardedCache(
	concurrency int,
	factory CacheFactory,
) *ShardedCache {
	var shards []Cache
	for i := 0; i < concurrency; i++ {
		shard := factory()
		shards = append(shards, shard)
	}

	return &ShardedCache{
		shards: shards,
	}
}

// Get ...
func (cache *ShardedCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	key := cache.makeKey(storage, color)
	return cache.shards[key].Get(
		storage,
		color,
	)
}

// Set ...
func (cache *ShardedCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	key := cache.makeKey(storage, color)
	cache.shards[key].Set(
		storage,
		color,
		move,
	)
}

func (cache *ShardedCache) makeKey(
	storage models.PieceStorage,
	color models.Color,
) int {
	return 0
}
