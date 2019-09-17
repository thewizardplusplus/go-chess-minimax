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

// IsSearchTerminated ...
func (
	group GroupTerminator,
) IsSearchTerminated(deep int) bool {
	terminators := group.terminators
	for _, terminator := range terminators {
		if terminator.IsSearchTerminated(deep) {
			return true
		}
	}

	return false
}
