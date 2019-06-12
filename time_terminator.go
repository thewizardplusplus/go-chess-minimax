package chessminimax

import (
	"time"
)

// TimeTerminator ...
type TimeTerminator struct {
	startTime       time.Time
	maximalDuration time.Duration
}

// NewTimeTerminator ...
func NewTimeTerminator(
	maximalDuration time.Duration,
) TimeTerminator {
	startTime := time.Now()
	return TimeTerminator{
		startTime:       startTime,
		maximalDuration: maximalDuration,
	}
}

// IsSearchTerminate ...
func (
	terminator TimeTerminator,
) IsSearchTerminate(deep int) bool {
	duration := time.Since(
		terminator.startTime,
	)
	return duration >=
		terminator.maximalDuration
}
