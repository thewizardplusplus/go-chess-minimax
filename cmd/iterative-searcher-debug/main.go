package main

import (
	"fmt"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func iterativeSearch(
	boardInFEN string,
	color models.Color,
	maximalDuration time.Duration,
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
	innerSearcher :=
		minimax.NewAlphaBetaSearcher(
			generator,
			// terminator will be set automatically
			// by the iterative searcher
			nil,
			evaluator,
		)

	terminator :=
		terminators.NewTimeTerminator(
			time.Now,
			maximalDuration,
		)
	searcher := minimax.NewIterativeSearcher(
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

func main() {
	move, err := iterativeSearch(
		"7K/8/8/8/8/8/pp6/kp6",
		models.Black,
		1000*time.Millisecond,
	)

	fmt.Println(move, err)
}
