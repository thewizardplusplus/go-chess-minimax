# go-chess-minimax

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-chess-minimax?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-chess-minimax)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-chess-minimax)](https://goreportcard.com/report/github.com/thewizardplusplus/go-chess-minimax)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-chess-minimax.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-chess-minimax)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-chess-minimax/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-chess-minimax)

The library that implements a chess engine based on the minimax algorithm.

_**Disclaimer:** this library was written directly on an Android smartphone with the AnGoIde IDE._

## Features

- move searcher used the [negamax algorithm](https://www.chessprogramming.org/Negamax);
- optimizations:
  - [alpha-beta pruning](https://www.chessprogramming.org/Alpha-Beta);
  - [transposition table](https://www.chessprogramming.org/Transposition_Table):
    - storing transpositions in an LRU cache;
    - hashing a transposition by its representation in [Forsyth–Edwards Notation](https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation);
    - replacing same transpositions on storing only if new one has a greater move quality:
      - move quality is directly proportional to a time of its evaluation;
    - sharing a [transposition table](https://www.chessprogramming.org/Transposition_Table) between searches;
    - [transposition table](https://www.chessprogramming.org/Transposition_Table) is safe for concurrent use (via a mutual exclusion lock over a whole storage);
  - [iterative deepening](https://www.chessprogramming.org/Iterative_Deepening);
  - parallel search ([Lazy SMP](https://www.chessprogramming.org/Iterative_Deepening)):
    - launch concurrent searches with same depths;
- searching termination:
  - by a deep;
  - by a time;
  - by calling a special method (it's safe for concurrent use);
- position evaluation only by a material (based on an [evaluation function](https://www.chessprogramming.org/Evaluation#Where_to_Start) of Claude Shannon);
- architecture features:
  - easily extensible and composable architecture of searching;
  - composable searching terminators.

## Installation

```
$ go get github.com/thewizardplusplus/go-chess-minimax
```

## Examples

`chessminimax.AlphaBetaSearcher.SearchMove()`:

```go
package main

import (
	"fmt"
	"log"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
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
```

`chessminimax.CachedSearcher.SearchMove()`:

```go
package main

import (
	"fmt"
	"log"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
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
```

`chessminimax.IterativeSearcher.SearchMove()`:

```go
package main

import (
	"fmt"
	"log"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
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
```

`chessminimax.ParallelSearcher.SearchMove()`:

```go
package main

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

func main() {
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
```

## Benchmarks

`chessminimax.AlphaBetaSearcher`:

```
BenchmarkAlphaBetaSearcher_1Ply-4   	     300	   4458950 ns/op
BenchmarkAlphaBetaSearcher_2Ply-4   	     100	  12019694 ns/op
BenchmarkAlphaBetaSearcher_3Ply-4   	      20	  74204644 ns/op
```

`chessminimax.CachedSearcher`:

```
BenchmarkCachedSearcher_1Ply-4      	   30000	     42505 ns/op
BenchmarkCachedSearcher_2Ply-4      	   30000	     44526 ns/op
BenchmarkCachedSearcher_3Ply-4      	   30000	     46041 ns/op
```

`chessminimax.IterativeSearcher`:

```
BenchmarkIterativeSearcher_1Ply-4   	     300	   4632475 ns/op
BenchmarkIterativeSearcher_2Ply-4   	    1000	   1443583 ns/op
BenchmarkIterativeSearcher_3Ply-4   	     500	   2319599 ns/op
```

`chessminimax.ParallelSearcher`:

```
BenchmarkParallelSearcher_1Ply-4    	     200	   6641428 ns/op
BenchmarkParallelSearcher_2Ply-4    	     500	   3607192 ns/op
BenchmarkParallelSearcher_3Ply-4    	     200	   5718786 ns/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus
