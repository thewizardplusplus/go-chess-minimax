package terminators

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminated(deep int) bool
}
