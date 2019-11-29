package chessminimax

import (
	"runtime"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func BenchmarkParallelSearcher_1Ply(benchmark *testing.B) {
	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	parallelCache := caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(parallelCache, initial, models.White, 1) // nolint: errcheck
	}
}

func BenchmarkParallelSearcher_2Ply(benchmark *testing.B) {
	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	parallelCache := caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(parallelCache, initial, models.White, 2) // nolint: errcheck
	}
}

func BenchmarkParallelSearcher_3Ply(benchmark *testing.B) {
	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	parallelCache := caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(parallelCache, initial, models.White, 3) // nolint: errcheck
	}
}

func parallelSearch(
	cache caches.Cache,
	boardInFEN string,
	color models.Color,
	maximalDeep int,
) (moves.ScoredMove, error) {
	storage, err :=
		uci.DecodePieceStorage(boardInFEN, pieces.NewPiece, models.NewBoard)
	if err != nil {
		return moves.ScoredMove{}, err
	}

	terminator := terminators.NewDeepTerminator(maximalDeep)
	searcher :=
		NewParallelSearcher(terminator, runtime.NumCPU(), func() MoveSearcher {
			var generator models.MoveGenerator
			var evaluator evaluators.MaterialEvaluator
			innerSearcher := NewAlphaBetaSearcher(
				generator,
				nil, // terminator will be set automatically by the iterative searcher
				evaluator,
			)

			// make and bind a cached searcher to inner one
			NewCachedSearcher(innerSearcher, cache)

			return NewIterativeSearcher(
				innerSearcher,
				nil, // terminator will be set automatically by the parallel searcher
			)
		})

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}
