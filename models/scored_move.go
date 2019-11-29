package models

import (
	"math"

	models "github.com/thewizardplusplus/go-chess-models"
)

// ScoredMove ...
type ScoredMove struct {
	Move    models.Move
	Score   float64
	Quality float64
}

// nolint: gochecknoglobals
var (
	initialScore = math.Inf(-1)
)

// NewScoredMove ...
func NewScoredMove() ScoredMove {
	return ScoredMove{Score: initialScore}
}

// IsUpdated ...
func (move ScoredMove) IsUpdated() bool {
	return move.Score != initialScore
}

// Update ...
func (move *ScoredMove) Update(
	scoredMove ScoredMove,
	rawMove models.Move,
	moveQuality float64,
) {
	score := -scoredMove.Score
	if move.Score >= score {
		return
	}

	*move = ScoredMove{
		Move:    rawMove,
		Score:   score,
		Quality: moveQuality,
	}
}
