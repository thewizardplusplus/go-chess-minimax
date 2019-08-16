package chessminimax

import (
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func BenchmarkCachedSearcher_1Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(
			cache,
			initial,
			models.White,
			1,
		)
	}
}

func BenchmarkCachedSearcher_2Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(
			cache,
			initial,
			models.White,
			2,
		)
	}
}

func BenchmarkCachedSearcher_3Ply(
	benchmark *testing.B,
) {
	cache := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(
			cache,
			initial,
			models.White,
			3,
		)
	}
}

func cachedSearch(
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
	terminator :=
		terminators.NewDeepTerminator(
			maximalDeep,
		)
	evaluator :=
		evaluators.MaterialEvaluator{}
	innerSearcher := NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)

	searcher := NewCachedSearcher(
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
