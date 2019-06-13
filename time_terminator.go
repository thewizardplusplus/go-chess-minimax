package chessminimax

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
	startTime := clock()
	return TimeTerminator{
		clock:           clock,
		maximalDuration: maximalDuration,
		startTime:       startTime,
	}
}

// IsSearchTerminate ...
func (
	terminator TimeTerminator,
) IsSearchTerminate(deep int) bool {
	currentTime := terminator.clock()
	duration := currentTime.Sub(
		terminator.startTime,
	)
	return duration >=
		terminator.maximalDuration
}
