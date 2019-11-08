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

// SearchProgress ...
func (
	group GroupTerminator,
) SearchProgress(deep int) float64 {
	if group.IsSearchTerminated(deep) {
		return 1
	}

	var maximalProgress float64
	terminators := group.terminators
	for _, terminator := range terminators {
		progress :=
			terminator.SearchProgress(deep)
		if progress > maximalProgress {
			maximalProgress = progress
		}
	}

	return maximalProgress
}
