package caches

import (
	"hash/fnv"
	"io"
	"strconv"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// CacheFactory ...
type CacheFactory func() Cache

// ShardedCache ...
type ShardedCache struct {
	shards   []Cache
	stringer Stringer
}

// NewShardedCache ...
func NewShardedCache(
	concurrency int,
	factory CacheFactory,
	stringer Stringer,
) *ShardedCache {
	var shards []Cache
	for i := 0; i < concurrency; i++ {
		shard := factory()
		shards = append(shards, shard)
	}

	return &ShardedCache{
		shards:   shards,
		stringer: stringer,
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
	text := cache.stringer(storage)
	text += " " + strconv.Itoa(int(color))

	hasher := fnv.New32()
	io.WriteString(hasher, text)

	hash := int(hasher.Sum32())
	hash %= len(cache.shards)

	return hash
}
