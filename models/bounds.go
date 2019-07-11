package models

import (
	"math"

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

// Next ...
func (bounds Bounds) Next() Bounds {
	alpha, beta := bounds.Alpha, bounds.Beta
	return Bounds{-beta, -alpha}
}

// Update ...
func (bounds *Bounds) Update(
	scoredMove ScoredMove,
	rawMove models.Move,
) (newScoredMove ScoredMove, ok bool) {
	score := -scoredMove.Score
	if score > bounds.Alpha {
		bounds.Alpha = score
	}
	if score >= bounds.Beta {
		return ScoredMove{rawMove, score}, false
	}

	return scoredMove, true
}
