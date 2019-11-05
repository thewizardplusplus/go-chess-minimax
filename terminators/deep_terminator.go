package terminators

// DeepTerminator ...
type DeepTerminator struct {
	maximalDeep int
}

// NewDeepTerminator ...
func NewDeepTerminator(
	maximalDeep int,
) DeepTerminator {
	return DeepTerminator{maximalDeep}
}

// IsSearchTerminated ...
func (
	terminator DeepTerminator,
) IsSearchTerminated(deep int) bool {
	return deep >= terminator.maximalDeep
}

// SearchProgress ...
func (
	terminator DeepTerminator,
) SearchProgress(deep int) float64 {
	return float64(deep) /
		float64(terminator.maximalDeep)
}
