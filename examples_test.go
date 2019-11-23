package chessminimax_test

import (
	"fmt"
	"log"
	"runtime"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func ExampleAlphaBetaSearcher() {
	storage, err :=
		uci.DecodePieceStorage("7K/8/7q/8/8/8/8/k7", pieces.NewPiece, models.NewBoard)
	if err != nil {
		log.Fatal(err)
	}

	var generator models.MoveGenerator
	var evaluator evaluators.MaterialEvaluator
	terminator := terminators.NewDeepTerminator(1)
	searcher := minimax.NewAlphaBetaSearcher(generator, terminator, evaluator)

	scoredMove, err :=
		searcher.SearchMove(storage, models.White, 0, moves.NewBounds())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", scoredMove)

	// Output: {Move:{Start:{File:7 Rank:7} Finish:{File:6 Rank:7}} Score:-9 Quality:1}
}

func ExampleCachedSearcher() {
	storage, err :=
		uci.DecodePieceStorage("7K/8/7q/8/8/8/8/k7", pieces.NewPiece, models.NewBoard)
	if err != nil {
		log.Fatal(err)
	}

	var generator models.MoveGenerator
	var evaluator evaluators.MaterialEvaluator
	terminator := terminators.NewDeepTerminator(1)
	innerSearcher := minimax.NewAlphaBetaSearcher(generator, terminator, evaluator)

	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	searcher := minimax.NewCachedSearcher(innerSearcher, cache)

	scoredMove, err :=
		searcher.SearchMove(storage, models.White, 0, moves.NewBounds())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", scoredMove)

	// Output: {Move:{Start:{File:7 Rank:7} Finish:{File:6 Rank:7}} Score:-9 Quality:1}
}

func ExampleIterativeSearcher() {
	storage, err :=
		uci.DecodePieceStorage("7K/8/7q/8/8/8/8/k7", pieces.NewPiece, models.NewBoard)
	if err != nil {
		log.Fatal(err)
	}

	var generator models.MoveGenerator
	var evaluator evaluators.MaterialEvaluator
	innerSearcher := minimax.NewAlphaBetaSearcher(
		generator,
		nil, // terminator will be set automatically by the iterative searcher
		evaluator,
	)

	// make and bind a cached searcher to inner one
	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	minimax.NewCachedSearcher(innerSearcher, cache)

	terminator := terminators.NewDeepTerminator(1)
	searcher := minimax.NewIterativeSearcher(innerSearcher, terminator)

	scoredMove, err :=
		searcher.SearchMove(storage, models.White, 0, moves.NewBounds())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", scoredMove)

	// Output: {Move:{Start:{File:7 Rank:7} Finish:{File:6 Rank:7}} Score:-9 Quality:1}
}

func ExampleParallelSearcher() {
	storage, err :=
		uci.DecodePieceStorage("7K/8/7q/8/8/8/8/k7", pieces.NewPiece, models.NewBoard)
	if err != nil {
		log.Fatal(err)
	}

	cache := caches.NewStringHashingCache(1e6, uci.EncodePieceStorage)
	terminator := terminators.NewDeepTerminator(1)
	searcher := minimax.NewParallelSearcher(
		terminator,
		runtime.NumCPU(),
		func() minimax.MoveSearcher {
			var generator models.MoveGenerator
			var evaluator evaluators.MaterialEvaluator
			innerSearcher := minimax.NewAlphaBetaSearcher(
				generator,
				nil, // terminator will be set automatically by the iterative searcher
				evaluator,
			)

			// make and bind a cached searcher to inner one
			minimax.NewCachedSearcher(innerSearcher, cache)

			return minimax.NewIterativeSearcher(
				innerSearcher,
				nil, // terminator will be set automatically by the parallel searcher
			)
		},
	)

	scoredMove, err :=
		searcher.SearchMove(storage, models.White, 0, moves.NewBounds())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", scoredMove)

	// Output: {Move:{Start:{File:7 Rank:7} Finish:{File:6 Rank:7}} Score:-9 Quality:1}
}
