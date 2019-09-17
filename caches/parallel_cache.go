package caches

import (
	"sync"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// ParallelCache ...
type ParallelCache struct {
	innerCache Cache
	locker     sync.RWMutex
}

// NewParallelCache ...
func NewParallelCache(
	innerCache Cache,
) *ParallelCache {
	return &ParallelCache{
		innerCache: innerCache,
	}
}

// Get ...
func (cache ParallelCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	cache.locker.RLock()
	defer cache.locker.RUnlock()

	return cache.innerCache.Get(
		storage,
		color,
	)
}

// Set ...
func (cache ParallelCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	cache.locker.Lock()
	defer cache.locker.Unlock()

	cache.innerCache.Set(
		storage,
		color,
		move,
	)
}
