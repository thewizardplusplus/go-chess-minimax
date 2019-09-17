package terminators

import (
	"sync/atomic"
)

// ParallelTerminator ...
type ParallelTerminator struct {
	terminationFlag uint64
}

// IsSearchTerminated ...
func (
	terminator *ParallelTerminator,
) IsSearchTerminated(deep int) bool {
	return terminator.terminationFlag != 0
}

// Terminate ...
func (
	terminator *ParallelTerminator,
) Terminate() {
	flag := &terminator.terminationFlag
	atomic.StoreUint64(flag, 1)
}
