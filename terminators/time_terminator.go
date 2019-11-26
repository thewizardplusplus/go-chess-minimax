package terminators

import (
	"time"
)

// Clock ...
type Clock func() time.Time

// TimeTerminator ...
type TimeTerminator struct {
	clock           Clock
	maximalDuration time.Duration
	startTime       time.Time
}

// NewTimeTerminator ...
func NewTimeTerminator(
	clock Clock,
	maximalDuration time.Duration,
) TimeTerminator {
	return TimeTerminator{
		clock:           clock,
		maximalDuration: maximalDuration,
		startTime:       clock(),
	}
}

// IsSearchTerminated ...
func (terminator TimeTerminator) IsSearchTerminated(deep int) bool {
	return terminator.elapsedTime() >= terminator.maximalDuration
}

// SearchProgress ...
func (terminator TimeTerminator) SearchProgress(deep int) float64 {
	if terminator.IsSearchTerminated(deep) {
		return 1
	}

	return float64(terminator.elapsedTime()) / float64(terminator.maximalDuration)
}

func (terminator TimeTerminator) elapsedTime() time.Duration {
	currentTime := terminator.clock()
	return currentTime.Sub(terminator.startTime)
}
