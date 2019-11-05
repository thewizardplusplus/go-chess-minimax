package terminators

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminated(deep int) bool
	SearchProgress(deep int) float64
}
