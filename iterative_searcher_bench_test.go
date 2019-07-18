package chessminimax

import (
  "testing"
  "time"

  "github.com/thewizardplusplus/go-chess-minimax/caches"
  "github.com/thewizardplusplus/go-chess-minimax/evaluators"
  moves "github.com/thewizardplusplus/go-chess-minimax/models"
  "github.com/thewizardplusplus/go-chess-minimax/terminators"
  models "github.com/thewizardplusplus/go-chess-models"
  "github.com/thewizardplusplus/go-chess-models/pieces"
)

func BenchmarkIterativeSearcher_1Ply(
  benchmark *testing.B,
) {
  for i := 0; i < benchmark.N; i++ {
    iterativeSearch(initial, models.White, 100*time.Millisecond)
  }
}

func BenchmarkIterativeSearcher_2Ply(
  benchmark *testing.B,
) {
  for i := 0; i < benchmark.N; i++ {
    iterativeSearch(initial, models.White, 200*time.Millisecond)
  }
}

func BenchmarkIterativeSearcher_3Ply(
  benchmark *testing.B,
) {
  for i := 0; i < benchmark.N; i++ {
    iterativeSearch(initial, models.White, 300*time.Millisecond)
  }
}

func iterativeSearch(
  boardInFEN string,
  color models.Color,
  maximalDuration time.Duration,
) (moves.ScoredMove, error) {
  storage, err := models.ParseBoard(
    boardInFEN,
    pieces.NewPiece,
  )
  if err != nil {
    return moves.ScoredMove{}, err
  }

  cache := make(caches.FENHashingCache)
  generator := models.MoveGenerator{}
  evaluator :=
    evaluators.MaterialEvaluator{}
  innerSearcher := NewAlphaBetaSearcher(
    generator,
    nil,
    evaluator,
  )
  NewCachedSearcher(cache, innerSearcher)

  terminator :=
    terminators.NewTimeTerminator(
      time.Now,
      maximalDuration,
    )
  searcher := NewIterativeSearcher(
    terminator,
    innerSearcher,
  )
  initialDeep := 0
  initialBounds := moves.NewBounds()
  return searcher.SearchMove(
    storage,
    color,
    initialDeep,
    initialBounds,
  )
}
