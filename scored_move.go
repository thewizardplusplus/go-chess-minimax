package chessminimax

import (
	"math"

	models "github.com/thewizardplusplus/go-chess-models"
)

// ScoredMove ...
type ScoredMove struct {
	Move  models.Move
	Score float64
}

var (
	initialScore = math.Inf(-1)
)

func newScoredMove() ScoredMove {
	return ScoredMove{Score: initialScore}
}

func (move ScoredMove) isUpdated() bool {
	return move.Score != initialScore
}

func (move *ScoredMove) update(
	scoredMove ScoredMove,
	rawMove models.Move,
) {
	score := -scoredMove.Score
	if move.Score < score {
		*move = ScoredMove{rawMove, score}
	}
}
