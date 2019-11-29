package chessminimax

import (
	"testing"

	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

var (
	initial = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

func BenchmarkAlphaBetaSearcher_1Ply(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 1) // nolint: errcheck
	}
}

func BenchmarkAlphaBetaSearcher_2Ply(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 2) // nolint: errcheck
	}
}

func BenchmarkAlphaBetaSearcher_3Ply(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		alphaBetaSearch(initial, models.White, 3) // nolint: errcheck
	}
}

func alphaBetaSearch(
	boardInFEN string,
	color models.Color,
	maximalDeep int,
) (moves.ScoredMove, error) {
	storage, err :=
		uci.DecodePieceStorage(boardInFEN, pieces.NewPiece, models.NewBoard)
	if err != nil {
		return moves.ScoredMove{}, err
	}

	var generator models.MoveGenerator
	var evaluator evaluators.MaterialEvaluator
	terminator := terminators.NewDeepTerminator(maximalDeep)
	searcher := NewAlphaBetaSearcher(generator, terminator, evaluator)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}
