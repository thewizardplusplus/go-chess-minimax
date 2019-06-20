package terminators

// SearchTerminator ...
type SearchTerminator interface {
	IsSearchTerminate(deep int) bool
}
