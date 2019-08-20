package caches

import (
	"container/list"

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

type bucket struct {
	key  key
	move moves.FailedMove
}

type bucketGroup map[key]*list.Element

// StringHashingCache ...
//
// It implements an LRU cache.
type StringHashingCache struct {
	buckets     bucketGroup
	queue       *list.List
	maximalSize int
	stringer    Stringer
}

// NewStringHashingCache ...
func NewStringHashingCache(
	maximalSize int,
	stringer Stringer,
) StringHashingCache {
	return StringHashingCache{
		buckets:     make(bucketGroup),
		queue:       list.New(),
		maximalSize: maximalSize,
		stringer:    stringer,
	}
}

// Get ...
func (cache StringHashingCache) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	key := cache.makeKey(storage, color)
	element, ok := cache.getElement(key)
	if !ok {
		return moves.FailedMove{}, false
	}

	return element.Value.(bucket).move, true
}

// Set ...
func (cache StringHashingCache) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	key := cache.makeKey(storage, color)
	newBucket := bucket{key, move}
	element, ok := cache.getElement(key)
	if ok {
		element.Value = newBucket
		return
	}

	element = cache.queue.PushFront(newBucket)
	cache.buckets[key] = element
	if cache.queue.Len() <=
		cache.maximalSize {
		return
	}

	element = cache.queue.Back()
	oldBucket := cache.queue.
		Remove(element).(bucket)
	delete(cache.buckets, oldBucket.key)
}

func (cache StringHashingCache) makeKey(
	storage models.PieceStorage,
	color models.Color,
) key {
	text := cache.stringer(storage)
	return key{text, color}
}

func (cache StringHashingCache) getElement(
	key key,
) (element *list.Element, ok bool) {
	element, ok = cache.buckets[key]
	if ok {
		cache.queue.MoveToFront(element)
	}

	return element, ok
}
