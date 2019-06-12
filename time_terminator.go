package chessminimax

import (
	"time"
)

// TimeTerminator ...
type TimeTerminator struct {
	clock           func() time.Time
	maximalDuration time.Duration
	startTime       time.Time
}

// NewTimeTerminator ...
func NewTimeTerminator(
	clock func() time.Time,
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
