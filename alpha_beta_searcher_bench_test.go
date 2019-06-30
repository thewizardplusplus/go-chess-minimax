package chessminimax

import (
	"math"
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func BenchmarkAlphaBetaSearcher_1Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 1)
	}
}

func BenchmarkAlphaBetaSearcher_2Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 2)
	}
}

func BenchmarkAlphaBetaSearcher_3Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 3)
	}
}

func alphaBetaSearch(
	boardInFEN string,
	color models.Color,
	maximalDeep int,
) (ScoredMove, error) {
	storage, err := models.ParseBoard(
		boardInFEN,
		pieces.NewPiece,
	)
	if err != nil {
		return ScoredMove{}, err
	}

	generator := models.MoveGenerator{}
	terminator :=
		terminators.NewDeepTerminator(
			maximalDeep,
		)
	evaluator :=
		evaluators.MaterialEvaluator{}
	searcher := NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)
	return searcher.SearchMove(
		storage,
		color,
		0,
		math.Inf(-1),
		math.Inf(+1),
	)
}
