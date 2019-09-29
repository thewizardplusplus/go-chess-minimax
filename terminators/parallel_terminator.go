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
	flag := &terminator.terminationFlag
	return atomic.LoadUint64(flag) != 0
}

// Terminate ...
func (
	terminator *ParallelTerminator,
) Terminate() {
	flag := &terminator.terminationFlag
	atomic.StoreUint64(flag, 1)
}
