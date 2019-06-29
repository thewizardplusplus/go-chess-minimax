package chessminimax

import (
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

var (
	initial = "rnbqkbnr/pppppppp/8/8" +
		"/8/8/PPPPPPPP/RNBQKBNR"
)

func BenchmarkNegamaxSearcher_1Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		negamaxSearch(initial, models.White, 1)
	}
}

func BenchmarkNegamaxSearcher_2Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		negamaxSearch(initial, models.White, 2)
	}
}

func BenchmarkNegamaxSearcher_3Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		negamaxSearch(initial, models.White, 3)
	}
}

func negamaxSearch(
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
	searcher := NewNegamaxSearcher(
		generator,
		terminator,
		evaluator,
	)
	return searcher.SearchMove(
		storage,
		color,
		0,
	)
}
