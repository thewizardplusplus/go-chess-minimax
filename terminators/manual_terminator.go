package terminators

import (
	"sync/atomic"
)

// ManualTerminator ...
type ManualTerminator struct {
	terminationFlag uint64
}

// IsSearchTerminated ...
func (terminator *ManualTerminator) IsSearchTerminated(deep int) bool {
	flag := &terminator.terminationFlag
	return atomic.LoadUint64(flag) != 0
}

// SearchProgress ...
func (terminator *ManualTerminator) SearchProgress(deep int) float64 {
	if terminator.IsSearchTerminated(deep) {
		return 1
	}

	return 0
}

// Terminate ...
func (terminator *ManualTerminator) Terminate() {
	flag := &terminator.terminationFlag
	atomic.StoreUint64(flag, 1)
}
