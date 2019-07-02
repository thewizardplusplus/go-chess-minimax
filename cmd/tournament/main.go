package main

import (
	"fmt"
	"log"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

const (
	gameCount       = 1
	maximalDeep     = 2
	maximalDuration = time.Second
	boardInFEN      = "rnbqk/ppppp/5" +
		"/PPPPP/RNBQK"
)

var (
	generator = models.MoveGenerator{}
	evaluator = evaluators.MaterialEvaluator{}
)

func makeTerminator(
	maximalDeep int,
	maximalDuration time.Duration,
) terminators.SearchTerminator {
	return terminators.NewGroupTerminator(
		terminators.NewDeepTerminator(
			maximalDeep,
		),
		terminators.NewTimeTerminator(
			time.Now,
			maximalDuration,
		),
	)
}

func negamaxSearch(
	storage models.PieceStorage,
	color models.Color,
	maximalDeep int,
	maximalDuration time.Duration,
) (minimax.ScoredMove, error) {
	terminator := makeTerminator(
		maximalDeep,
		maximalDuration,
	)
	searcher := minimax.NewNegamaxSearcher(
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

func alphaBetaSearch(
	storage models.PieceStorage,
	color models.Color,
	maximalDeep int,
	maximalDuration time.Duration,
) (minimax.ScoredMove, error) {
	terminator := makeTerminator(
		maximalDeep,
		maximalDuration,
	)
	bounds := minimax.NewBounds()
	searcher := minimax.NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)
	return searcher.SearchMove(
		storage,
		color,
		0,
		bounds,
	)
}

func game(
	storage models.PieceStorage,
	color models.Color,
	maximalDeep int,
	maximalDuration time.Duration,
) (models.Color, error) {
	for {
		move, err := negamaxSearch(
			storage,
			color,
			maximalDeep,
			maximalDuration,
		)
		if err != nil {
			return color, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()

		move, err = alphaBetaSearch(
			storage,
			color,
			maximalDeep,
			maximalDuration,
		)
		if err != nil {
			return color, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()
	}
}

func main() {
	storage, err := models.ParseBoard(
		boardInFEN,
		pieces.NewPiece,
	)
	if err != nil {
		log.Fatal(err)
	}

	scores := make(map[models.Color]float64)
	color := models.White
	for i := 0; i < gameCount; i++ {
		loserColor, err := game(
			storage,
			color,
			maximalDeep,
			maximalDuration,
		)
		switch err {
		case minimax.ErrCheckmate:
			scores[loserColor.Negative()]++
		case minimax.ErrDraw:
			scores[models.Black] += 0.5
			scores[models.White] += 0.5
		}

		color = color.Negative()
	}

	fmt.Println(scores)
}
