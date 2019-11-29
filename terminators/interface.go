package terminators

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminated(deep int) bool

	// It should return a value between 0 and 1 inclusive.
	SearchProgress(deep int) float64
}
