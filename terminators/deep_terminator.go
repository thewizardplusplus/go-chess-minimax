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

// IsSearchTerminate ...
func (
	terminator DeepTerminator,
) IsSearchTerminate(deep int) bool {
	return deep >= terminator.maximalDeep
}
