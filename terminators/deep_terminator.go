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
