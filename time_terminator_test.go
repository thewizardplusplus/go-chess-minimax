package chessminimax

import (
	"testing"
	"time"
)

func clock() time.Time {
	year, month, day := 2006, time.January, 2
	hour, minute, second := 15, 4, 5
	return time.Date(
		year, month, day,
		hour, minute, second,
		int(0*time.Nanosecond),
		time.UTC,
	)
}

func TestNewTimeTerminator(test *testing.T) {
	terminator := NewTimeTerminator(
		clock,
		5*time.Second,
	)

	if terminator.clock == nil {
		test.Fail()
	}
	if terminator.maximalDuration !=
		5*time.Second {
		test.Fail()
	}
	if !terminator.startTime.
		Equal(clock()) {
		test.Fail()
	}
}
