package terminators

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimeTerminator(test *testing.T) {
	var wrappedClockCallCount int
	wrappedClock := func() time.Time { wrappedClockCallCount++; return clock() }
	maximalDuration := 5 * time.Second
	terminator := NewTimeTerminator(wrappedClock, maximalDuration)

	gotClock := reflect.ValueOf(terminator.clock).Pointer()
	wantClock := reflect.ValueOf(wrappedClock).Pointer()
	if gotClock != wantClock {
		test.Fail()
	}

	if terminator.maximalDuration != maximalDuration {
		test.Fail()
	}

	startTime := clock()
	if !terminator.startTime.Equal(startTime) {
		test.Fail()
	}

	if wrappedClockCallCount != 1 {
		test.Fail()
	}
}

func TestTimeTerminatorIsSearchTerminated(test *testing.T) {
	type fields struct {
		clock           Clock
		maximalDuration time.Duration
		startTime       time.Time
	}
	type args struct {
		deep int
	}
	type data struct {
		fields fields
		args   args
		want   bool
	}

	var clockCallCount int
	for _, data := range []data{
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-4 * time.Second),
			},
			args: args{5},
			want: false,
		},
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-5 * time.Second),
			},
			args: args{5},
			want: true,
		},
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-6 * time.Second),
			},
			args: args{5},
			want: true,
		},
	} {
		clockCallCount = 0

		terminator := TimeTerminator{
			clock:           data.fields.clock,
			maximalDuration: data.fields.maximalDuration,
			startTime:       data.fields.startTime,
		}
		got := terminator.IsSearchTerminated(data.args.deep)

		if got != data.want {
			test.Fail()
		}
		if clockCallCount != 1 {
			test.Fail()
		}
	}
}

func TestTimeTerminatorSearchProgress(test *testing.T) {
	type fields struct {
		clock           Clock
		maximalDuration time.Duration
		startTime       time.Time
	}
	type args struct {
		deep int
	}
	type data struct {
		fields fields
		args   args
		want   float64
	}

	var clockCallCount int
	for _, data := range []data{
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 100 * time.Second,
				startTime:       clock(),
			},
			args: args{5},
			want: 0,
		},
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-75 * time.Second),
			},
			args: args{5},
			want: 0.75,
		},
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-100 * time.Second),
			},
			args: args{5},
			want: 1,
		},
		{
			fields: fields{
				clock:           func() time.Time { clockCallCount++; return clock() },
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-110 * time.Second),
			},
			args: args{5},
			want: 1,
		},
	} {
		clockCallCount = 0

		terminator := TimeTerminator{
			clock:           data.fields.clock,
			maximalDuration: data.fields.maximalDuration,
			startTime:       data.fields.startTime,
		}
		got := terminator.SearchProgress(data.args.deep)

		if got != data.want {
			test.Fail()
		}
		if clockCallCount != 1 {
			test.Fail()
		}
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
