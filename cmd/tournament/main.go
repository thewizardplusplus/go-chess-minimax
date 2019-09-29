package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

// TaskInbox ...
type TaskInbox chan func()

// Side ...
type Side int

// ...
const (
	Iterative Side = iota
	Parallel
)

// Score ...
type Score struct {
	winCount int64
	score    int64
}

// Win ...
func (score *Score) Win() {
	atomic.AddInt64(&score.winCount, 1)
	atomic.AddInt64(&score.score, 10)
}

// Draw ...
func (score *Score) Draw() {
	atomic.AddInt64(&score.score, 5)
}

// Score ...
func (score Score) Score() float64 {
	return float64(score.score) / 10
}

// Elo ...
func (score Score) Elo(
	gameCount int64,
) float64 {
	winPercent := float64(score.winCount) /
		float64(gameCount)
	return 400 *
		math.Log10(winPercent/(1-winPercent))
}

// Scores ...
type Scores struct {
	gameCount int64
	iterative Score
	parallel  Score
}

// AddGame ...
func (scores *Scores) AddGame(
	loserSide Side,
	err error,
) {
	atomic.AddInt64(&scores.gameCount, 1)

	switch err {
	case minimax.ErrCheckmate:
		switch loserSide {
		case Iterative:
			scores.parallel.Win()
		case Parallel:
			scores.iterative.Win()
		}
	case minimax.ErrDraw:
		scores.iterative.Draw()
		scores.parallel.Draw()
	}
}

// String ...
func (scores Scores) String() string {
	return fmt.Sprintf(
		"Games: %d "+
			"Iterative: %.1f "+
			"Parallel: %.1f\n"+
			"Parallel Elo Delta: %.2f",
		scores.gameCount,
		scores.iterative.Score(),
		scores.parallel.Score(),
		scores.parallel.Elo(scores.gameCount),
	)
}

const (
	gameCount    = 10
	maxDeep      = 4
	maxDuration  = 500 * time.Millisecond
	maxMoveCount = 40
	boardInFEN   = "rnbqk/ppppp/5/PPPPP/RNBQK"
)

var (
	generator = models.MoveGenerator{}
	evaluator = evaluators.MaterialEvaluator{}

	errTooLong = errors.New("too long")
)

func pool() (tasks TaskInbox, wait func()) {
	threadCount := runtime.NumCPU()

	var waiter sync.WaitGroup
	waiter.Add(threadCount)

	tasks = make(TaskInbox)
	for i := 0; i < threadCount; i++ {
		go func() {
			defer waiter.Done()
			fmt.Print("#")

			for task := range tasks {
				fmt.Print("%")
				task()
			}
		}()
	}

	return tasks, func() { waiter.Wait() }
}

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

func iterativeSearch(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (moves.ScoredMove, error) {
	innerSearcher :=
		minimax.NewAlphaBetaSearcher(
			generator,
			// terminator will be set
			// automatically
			// by the iterative searcher
			nil,
			evaluator,
		)

	// make and bind a cached searcher
	// to inner one
	minimax.NewCachedSearcher(
		innerSearcher,
		cache,
	)

	terminator := makeTerminator(
		maxDeep,
		maxDuration,
	)
	searcher := minimax.NewIterativeSearcher(
		innerSearcher,
		terminator,
	)

	return searcher.SearchMove(
		storage,
		color,
		0,
		moves.NewBounds(),
	)
}

func parallelSearch(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (moves.ScoredMove, error) {
	baseTerminator := makeTerminator(
		maxDeep,
		maxDuration,
	)

	searcher := minimax.NewParallelSearcher(
		runtime.NumCPU(),
		func(
			parallelTerminator terminators.SearchTerminator,
		) minimax.MoveSearcher {
			innerSearcher :=
				minimax.NewAlphaBetaSearcher(
					generator,
					// terminator will be set
					// automatically
					// by the Parallel searcher
					nil,
					evaluator,
				)

			// make and bind a cached searcher
			// to inner one
			minimax.NewCachedSearcher(
				innerSearcher,
				cache,
			)

			terminator :=
				terminators.NewGroupTerminator(
					baseTerminator,
					parallelTerminator,
				)
			return minimax.NewIterativeSearcher(
				innerSearcher,
				terminator,
			)
		},
	)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}

func game(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
	maxMoveCount int,
) (Side, error) {
	cacheOne := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	cacheTwo := caches.NewStringHashingCache(
		1e6,
		uci.EncodePieceStorage,
	)
	parallelCacheTwo :=
		caches.NewParallelCache(cacheTwo)

	for i := 0; i < maxMoveCount; i++ {
		if i%5 == 0 {
			fmt.Print(".")
		}

		move, err := iterativeSearch(
			cacheOne,
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return Iterative, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()

		move, err = parallelSearch(
			parallelCacheTwo,
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return Parallel, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()
	}

	return 0, errTooLong
}

func markGame(loserSide Side, err error) {
	switch err {
	case minimax.ErrCheckmate:
		switch loserSide {
		case Iterative:
			fmt.Print("P")
		case Parallel:
			fmt.Print("I")
		}
	case minimax.ErrDraw:
		fmt.Print("D")
	case errTooLong:
		fmt.Print("L")
	}
}

func main() {
	start := time.Now()
	storage, err := uci.DecodePieceStorage(
		boardInFEN,
		pieces.NewPiece,
		models.NewBoard,
	)
	if err != nil {
		log.Fatal(err)
	}

	var scores Scores
	tasks, wait := pool()
	initialColor := models.White
	for scores.gameCount < gameCount {
		initialColorCopy := initialColor
		tasks <- func() {
			loserSide, err := game(
				storage,
				initialColorCopy,
				maxDeep,
				maxDuration,
				maxMoveCount,
			)
			markGame(loserSide, err)
			if err != errTooLong {
				scores.AddGame(loserSide, err)
			}
		}

		initialColor = initialColor.Negative()
	}

	close(tasks)
	wait()

	fmt.Println()
	fmt.Println(scores)
	fmt.Println(time.Since(start))
}
