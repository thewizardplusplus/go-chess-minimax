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
) ShardedCache {
	var shards []Cache
	for i := 0; i < concurrency; i++ {
		shard := factory()
		shards = append(shards, shard)
	}

	return ShardedCache{
		shards:   shards,
		stringer: stringer,
	}
}

// Get ...
func (cache ShardedCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	index := cache.makeIndex(storage, color)
	return cache.shards[index].Get(
		storage,
		color,
	)
}

// Set ...
func (cache ShardedCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	index := cache.makeIndex(storage, color)
	cache.shards[index].Set(
		storage,
		color,
		move,
	)
}

func (cache ShardedCache) makeIndex(
	storage models.PieceStorage,
	color models.Color,
) int {
	text := cache.stringer(storage)
	text += " " + strconv.Itoa(int(color))

	hasher := fnv.New32()
	io.WriteString(hasher, text)

	hash := hasher.Sum32()
	hash %= uint32(len(cache.shards))

	return int(hash)
}
