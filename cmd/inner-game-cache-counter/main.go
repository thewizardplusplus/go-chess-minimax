package main

import (
	"fmt"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

// CacheWrapper ...
type CacheWrapper struct {
	caches.Cache

	getCount int
	setCount int
}

// NewCacheWrapper ...
func NewCacheWrapper(
	cache caches.Cache,
) *CacheWrapper {
	return &CacheWrapper{Cache: cache}
}

// Get ...
func (cache *CacheWrapper) Get(
	storage models.PieceStorage,
	color models.Color,
) (move moves.FailedMove, ok bool) {
	move, ok = cache.Cache.Get(storage, color)
	if ok {
		cache.getCount++
	}

	return move, ok
}

// Set ...
func (cache *CacheWrapper) Set(
	storage models.PieceStorage,
	color models.Color,
	move moves.FailedMove,
) {
	_, ok := cache.Cache.Get(storage, color)
	if !ok {
		cache.setCount++
	}

	cache.Cache.Set(storage, color, move)
}

func cachedSearch(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	cache caches.Cache,
) (moves.ScoredMove, error) {
	generator := models.MoveGenerator{}
	terminator :=
		terminators.NewDeepTerminator(maxDeep)
	evaluator :=
		evaluators.MaterialEvaluator{}
	innerSearcher :=
		minimax.NewAlphaBetaSearcher(
			generator,
			terminator,
			evaluator,
		)

	searcher := minimax.NewCachedSearcher(
		innerSearcher,
		cache,
	)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}

func game(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxMoveCount int,
) error {
	start := time.Now()

	var ply int
	for ; ply < maxMoveCount*2; ply++ {
		cache := caches.NewStringHashingCache(
			uci.EncodePieceStorage,
		)
		wrappedCache := NewCacheWrapper(cache)

		move, err := cachedSearch(
			storage,
			color,
			maxDeep,
			wrappedCache,
		)
		if err != nil {
			return err
		}

		fmt.Printf(
			"ply: %d, count: %d/%d, time: %s\n",
			ply,
			wrappedCache.getCount,
			wrappedCache.setCount,
			time.Since(start),
		)

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()
	}

	return nil
}

func main() {
	type data struct {
		name         string
		fen          string
		maxDeep      int
		maxMoveCount int
	}

	for _, data := range []data{
		data{
			name: "usual chess",
			fen: "rnbqkbnr/pppppppp/8/8" +
				"/8/8/PPPPPPPP/RNBQKBNR",
			maxDeep:      4,
			maxMoveCount: 5,
		},
		data{
			name: "minichess",
			fen: "rnbqk/ppppp/5" +
				"/PPPPP/RNBQK",
			maxDeep:      4,
			maxMoveCount: 6,
		},
	} {
		fmt.Printf(
			"name: %s\nfen: %s\n",
			data.name,
			data.fen,
		)

		storage, err := uci.DecodePieceStorage(
			data.fen,
			pieces.NewPiece,
			models.NewBoard,
		)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		err = game(
			storage,
			models.White,
			data.maxDeep,
			data.maxMoveCount,
		)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		fmt.Println()
	}
}
