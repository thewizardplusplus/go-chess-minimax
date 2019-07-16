package chessminimax

import (
	"reflect"
	"testing"
	"time"
)

func TestNewIterativeSearcher(
	test *testing.T,
) {
	innerSearcher := MockMoveSearcher{
		setSearcher: func(
			innerSearcher MoveSearcher,
		) {
			mock, ok :=
				innerSearcher.(*IterativeSearcher)
			if !ok || mock == nil {
				test.Fail()
			}
		},
	}

	maximalDuration := 5 * time.Second
	searcher := NewIterativeSearcher(
		innerSearcher,
		clock,
		maximalDuration,
	)

	_, ok := searcher.
		MoveSearcher.(MockMoveSearcher)
	if !ok {
		test.Fail()
	}

	gotClock := reflect.
		ValueOf(searcher.clock).
		Pointer()
	wantClock := reflect.
		ValueOf(clock).
		Pointer()
	if gotClock != wantClock {
		test.Fail()
	}

	if searcher.maximalDuration !=
		maximalDuration {
		test.Fail()
	}
}

func clock() time.Time {
	year, month, day := 2006, time.January, 2
	hour, minute, second := 15, 4, 5
	return time.Date(
		year, month, day,
		hour, minute, second,
		0,        // nanosecond
		time.UTC, // location
	)
}
