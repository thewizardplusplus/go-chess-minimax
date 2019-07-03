package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

type ScoreGroup struct {
	locker    sync.Mutex
	gameCount int
	negamax   float64
	alphaBeta float64
}

func (scores *ScoreGroup) AddGame(
	initialColor models.Color,
	loserColor models.Color,
	err error,
) {
	scores.locker.Lock()
	defer scores.locker.Unlock()

	scores.gameCount++

	switch err {
	case minimax.ErrCheckmate:
		if loserColor != initialColor {
			scores.negamax++
		} else {
			scores.alphaBeta++
		}
	case minimax.ErrDraw:
		scores.negamax += 0.5
		scores.alphaBeta += 0.5
	}
}

func (scores ScoreGroup) String() string {
	return fmt.Sprintf(
		"Games: %d Negamax: %f Alpha-Beta: %f",
		scores.gameCount,
		scores.negamax,
		scores.alphaBeta,
	)
}

const (
	gameCount    = 1
	maxDeep      = 4
	maxDuration  = 500 * time.Millisecond
	maxMoveCount = 20
	boardInFEN   = "rnbqk/ppppp/5/PPPPP/RNBQK"
)

var (
	generator = models.MoveGenerator{}
	evaluator = evaluators.MaterialEvaluator{}

	errTooLong = errors.New("too long")
)

func makeTerminator(
	maxDeep int,
	maxDuration time.Duration,
) terminators.SearchTerminator {
	return terminators.NewGroupTerminator(
		terminators.NewDeepTerminator(
			maxDeep,
		),
		terminators.NewTimeTerminator(
			time.Now,
			maxDuration,
		),
	)
}

func negamaxSearch(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (minimax.ScoredMove, error) {
	terminator := makeTerminator(
		maxDeep,
		maxDuration,
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
	maxDeep int,
	maxDuration time.Duration,
) (minimax.ScoredMove, error) {
	terminator := makeTerminator(
		maxDeep,
		maxDuration,
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
	maxDeep int,
	maxDuration time.Duration,
	maxMoveCount int,
) (models.Color, error) {
	for i := 0; i < maxMoveCount; i++ {
		if i%5 == 0 {
			fmt.Print(".")
		}

		move, err := negamaxSearch(
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return color, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()

		move, err = alphaBetaSearch(
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return color, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()
	}

	return 0, errTooLong
}

func main() {
	start := time.Now()
	storage, err := models.ParseBoard(
		boardInFEN,
		pieces.NewPiece,
	)
	if err != nil {
		log.Fatal(err)
	}

	var scores ScoreGroup
	threadCount := runtime.NumCPU()
	initialColor := models.White
	for i := 0; i < gameCount; i++ {
		var waiter sync.WaitGroup
		waiter.Add(threadCount)

		for j := 0; j < threadCount; j++ {
			go func() {
				defer waiter.Done()

				loserColor, err := game(
					storage,
					initialColor,
					maxDeep,
					maxDuration,
					maxMoveCount,
				)
				if err == errTooLong {
					return
				}

				scores.AddGame(
					initialColor,
					loserColor,
					err,
				)
			}()
		}
		waiter.Wait()

		initialColor = initialColor.Negative()
	}

	fmt.Println()
	fmt.Println(scores)
	fmt.Println(time.Since(start))
}
