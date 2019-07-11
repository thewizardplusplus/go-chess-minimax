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
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

// TaskInbox ...
type TaskInbox chan func()

// Side ...
type Side int

// ...
const (
	AlphaBeta Side = iota
	Cached
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
	alphaBeta Score
	cached    Score
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
		case AlphaBeta:
			scores.cached.Win()
		case Cached:
			scores.alphaBeta.Win()
		}
	case minimax.ErrDraw:
		scores.alphaBeta.Draw()
		scores.cached.Draw()
	}
}

// String ...
func (scores Scores) String() string {
	return fmt.Sprintf(
		"Games: %d "+
			"Alpha-Beta: %.1f "+
			"Cached: %.1f\n"+
			"Cached Elo Delta: %.2f",
		scores.gameCount,
		scores.alphaBeta.Score(),
		scores.cached.Score(),
		scores.cached.Elo(scores.gameCount),
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

func alphaBetaSearch(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (moves.ScoredMove, error) {
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

func cachedSearch(
	storage models.PieceStorage,
	color models.Color,
	maxDeep int,
	maxDuration time.Duration,
) (moves.ScoredMove, error) {
	cache := make(minimax.FENHashingCache)
	terminator := makeTerminator(
		maxDeep,
		maxDuration,
	)
	bounds := minimax.NewBounds()
	innerSearcher :=
		minimax.NewAlphaBetaSearcher(
			generator,
			terminator,
			evaluator,
		)
	searcher := minimax.NewCachedSearcher(
		cache,
		innerSearcher,
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
) (Side, error) {
	for i := 0; i < maxMoveCount; i++ {
		if i%5 == 0 {
			fmt.Print(".")
		}

		move, err := alphaBetaSearch(
			storage,
			color,
			maxDeep,
			maxDuration,
		)
		if err != nil {
			return AlphaBeta, err
		}

		storage = storage.ApplyMove(move.Move)
		color = color.Negative()

		move, err = cachedSearch(
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
		case AlphaBeta:
			fmt.Print("C")
		case Cached:
			fmt.Print("A")
		}
	case minimax.ErrDraw:
		fmt.Print("D")
	case errTooLong:
		fmt.Print("L")
	}
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
