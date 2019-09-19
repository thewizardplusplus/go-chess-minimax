package chessminimax

import (
	"runtime"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func BenchmarkParallelSearcher_1Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	parallelCache :=
		caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(
			parallelCache,
			initial,
			models.White,
			1,
		)
	}
}

func BenchmarkParallelSearcher_2Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	parallelCache :=
		caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(
			parallelCache,
			initial,
			models.White,
			2,
		)
	}
}

func BenchmarkParallelSearcher_3Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	parallelCache :=
		caches.NewParallelCache(cache)
	for i := 0; i < benchmark.N; i++ {
		parallelSearch(
			parallelCache,
			initial,
			models.White,
			3,
		)
	}
}

func parallelSearch(
	cache caches.Cache,
	boardInFEN string,
	color models.Color,
	maximalDeep int,
) (moves.ScoredMove, error) {
	storage, err := uci.DecodePieceStorage(
		boardInFEN,
		pieces.NewPiece,
		models.NewBoard,
	)
	if err != nil {
		return moves.ScoredMove{}, err
	}

	generator := models.MoveGenerator{}
	baseTerminator :=
		terminators.NewDeepTerminator(
			maximalDeep,
		)
	evaluator :=
		evaluators.MaterialEvaluator{}

	searcher := NewParallelSearcher(
		runtime.NumCPU(),
		func(
			parallelTerminator terminators.SearchTerminator,
		) MoveSearcher {
			innerSearcher := NewAlphaBetaSearcher(
				generator,
				// terminator will be set
				// automatically
				// by the Parallel searcher
				nil,
				evaluator,
			)

			// make and bind a cached searcher
			// to inner one
			NewCachedSearcher(innerSearcher, cache)

			terminator :=
				terminators.NewGroupTerminator(
					baseTerminator,
					parallelTerminator,
				)
			return NewIterativeSearcher(
				innerSearcher,
				terminator,
			)
		},
	)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}
