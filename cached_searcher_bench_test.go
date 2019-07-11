package chessminimax

import (
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func BenchmarkCachedSearcher_1Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(initial, models.White, 1)
	}
}

func BenchmarkCachedSearcher_2Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(initial, models.White, 2)
	}
}

func BenchmarkCachedSearcher_3Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		cachedSearch(initial, models.White, 3)
	}
}

func cachedSearch(
	boardInFEN string,
	color models.Color,
	maximalDeep int,
) (moves.ScoredMove, error) {
	storage, err := models.ParseBoard(
		boardInFEN,
		pieces.NewPiece,
	)
	if err != nil {
		return moves.ScoredMove{}, err
	}

	cache := make(caches.FENHashingCache)
	generator := models.MoveGenerator{}
	terminator :=
		terminators.NewDeepTerminator(
			maximalDeep,
		)
	evaluator :=
		evaluators.MaterialEvaluator{}
	initialDeep := 0
	initialBounds := moves.NewBounds()
	innerSearcher := NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)
	searcher := NewCachedSearcher(
		cache,
		innerSearcher,
	)
	return searcher.SearchMove(
		storage,
		color,
		initialDeep,
		initialBounds,
	)
}
