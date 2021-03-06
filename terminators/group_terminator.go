package terminators

// GroupTerminator ...
type GroupTerminator struct {
	terminators []SearchTerminator
}

// NewGroupTerminator ...
func NewGroupTerminator(terminators ...SearchTerminator) GroupTerminator {
	return GroupTerminator{terminators}
}

// IsSearchTerminated ...
func (group GroupTerminator) IsSearchTerminated(deep int) bool {
	for _, terminator := range group.terminators {
		if terminator.IsSearchTerminated(deep) {
			return true
		}
	}

	return false
}

// SearchProgress ...
func (group GroupTerminator) SearchProgress(deep int) float64 {
	var maximalProgress float64
	for _, terminator := range group.terminators {
		progress := terminator.SearchProgress(deep)
		if progress > maximalProgress {
			maximalProgress = progress
		}
	}

	return maximalProgress
}
