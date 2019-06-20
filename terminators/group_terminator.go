package terminators

// GroupTerminator ...
type GroupTerminator struct {
	terminators []SearchTerminator
}

// NewGroupTerminator ...
func NewGroupTerminator(
	terminators ...SearchTerminator,
) GroupTerminator {
	return GroupTerminator{terminators}
}

// IsSearchTerminate ...
func (
	group GroupTerminator,
) IsSearchTerminate(deep int) bool {
	terminators := group.terminators
	for _, terminator := range terminators {
		if terminator.IsSearchTerminate(deep) {
			return true
		}
	}

	return false
}
