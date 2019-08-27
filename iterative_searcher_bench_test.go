package chessminimax

import (
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func BenchmarkIterativeSearcher_1Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		iterativeSearch(initial, models.White, 1)
	}
}

func BenchmarkIterativeSearcher_2Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		iterativeSearch(initial, models.White, 2)
	}
}

func BenchmarkIterativeSearcher_3Ply(
	benchmark *testing.B,
) {
	for i := 0; i < benchmark.N; i++ {
		iterativeSearch(initial, models.White, 3)
	}
}

func iterativeSearch(
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
	evaluator :=
		evaluators.MaterialEvaluator{}
	innerSearcher := NewAlphaBetaSearcher(
		generator,
		// terminator will be set automatically
		// by the iterative searcher
		nil,
		evaluator,
	)

	terminator :=
		terminators.NewDeepTerminator(
			maximalDeep,
		)
	searcher := NewIterativeSearcher(
		innerSearcher,
		terminator,
	)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}
