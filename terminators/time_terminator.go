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
	startTime := clock()
	return TimeTerminator{
		clock:           clock,
		maximalDuration: maximalDuration,
		startTime:       startTime,
	}
}

// IsSearchTerminated ...
func (terminator TimeTerminator) IsSearchTerminated(deep int) bool {
	_, ok := terminator.duration()
	return !ok
}

// SearchProgress ...
func (terminator TimeTerminator) SearchProgress(deep int) float64 {
	duration, ok := terminator.duration()
	if !ok {
		return 1
	}

	return float64(duration) / float64(terminator.maximalDuration)
}

func (terminator TimeTerminator) duration() (duration time.Duration, ok bool) {
	currentTime := terminator.clock()
	duration = currentTime.Sub(terminator.startTime)
	ok = duration < terminator.maximalDuration
	return duration, ok
}
