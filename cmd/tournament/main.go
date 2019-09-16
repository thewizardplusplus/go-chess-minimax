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
	Cached Side = iota
	Iterative
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
	cached    Score
	iterative Score
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
		case Cached:
			scores.iterative.Win()
		case Iterative:
			scores.cached.Win()
		}
	case minimax.ErrDraw:
		scores.cached.Draw()
		scores.iterative.Draw()
	}
}

// String ...
func (scores Scores) String() string {
	return fmt.Sprintf(
		"Games: %d "+
			"Cached: %.1f "+
			"Iterative: %.1f\n"+
			"Iterative Elo Delta: %.2f",
		scores.gameCount,
		scores.cached.Score(),
		scores.iterative.Score(),
		scores.iterative.Elo(scores.gameCount),
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

func cachedSearch(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (moves.ScoredMove, error) {
	terminator := makeTerminator(
		maxDeep,
		maxDuration,
	)
	bounds := moves.NewBounds()
	innerSearcher :=
		minimax.NewAlphaBetaSearcher(
			generator,
			terminator,
			evaluator,
		)
	searcher := minimax.NewCachedSearcher(
		innerSearcher,
		cache,
	)
	return searcher.SearchMove(
		storage,
		color,
		0,
		bounds,
	)
}

func iterativeSearch(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	maximalDeep int,
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

	for i := 0; i < maxMoveCount; i++ {
		if i%5 == 0 {
			fmt.Print(".")
		}

		move, err := cachedSearch(
			cacheOne,
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return Cached, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()

		move, err = iterativeSearch(
			cacheTwo,
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return Cached, err
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
		case Cached:
			fmt.Print("I")
		case Iterative:
			fmt.Print("C")
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
