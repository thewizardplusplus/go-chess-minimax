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
func (
	terminator TimeTerminator,
) IsSearchTerminated(deep int) bool {
	currentTime := terminator.clock()
	duration := currentTime.
		Sub(terminator.startTime)
	return duration >=
		terminator.maximalDuration
}
