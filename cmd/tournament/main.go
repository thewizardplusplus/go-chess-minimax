package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

// ScoreGroup ...
type ScoreGroup struct {
	gameCount int64
	negamax   int64
	alphaBeta int64
}

// AddGame ...
func (scores *ScoreGroup) AddGame(
	initialColor models.Color,
	loserColor models.Color,
	err error,
) {
	atomic.AddInt64(&scores.gameCount, 1)

	switch err {
	case minimax.ErrCheckmate:
		if loserColor != initialColor {
			atomic.AddInt64(&scores.negamax, 10)
		} else {
			atomic.AddInt64(&scores.alphaBeta, 10)
		}
	case minimax.ErrDraw:
		atomic.AddInt64(&scores.negamax, 5)
		atomic.AddInt64(&scores.alphaBeta, 5)
	}
}

// String ...
func (scores ScoreGroup) String() string {
	return fmt.Sprintf(
		"Games: %d Negamax: %f Alpha-Beta: %f",
		scores.gameCount,
		float64(scores.negamax)/10,
		float64(scores.alphaBeta)/10,
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
