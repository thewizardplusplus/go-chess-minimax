package chessminimax

// Bounds ...
type Bounds struct {
	Alpha float64
	Beta  float64
}

func (bounds Bounds) next() Bounds {
	alpha, beta := bounds.Alpha, bounds.Beta
	return Bounds{-beta, -alpha}
}

func (bounds *Bounds) update(
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
