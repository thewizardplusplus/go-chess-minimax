package chessminimax

import (
	"math"

	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// Bounds ...
type Bounds struct {
	Alpha float64
	Beta  float64
}

// NewBounds ...
func NewBounds() Bounds {
	return Bounds{math.Inf(-1), math.Inf(+1)}
}

func (bounds Bounds) next() Bounds {
	alpha, beta := bounds.Alpha, bounds.Beta
	return Bounds{-beta, -alpha}
}

func (bounds *Bounds) update(
	scoredMove moves.ScoredMove,
	rawMove models.Move,
) (newScoredMove moves.ScoredMove, ok bool) {
	score := -scoredMove.Score
	if score > bounds.Alpha {
		bounds.Alpha = score
	}
	if score >= bounds.Beta {
		return moves.ScoredMove{rawMove, score},
			false
	}

	return scoredMove, true
}
