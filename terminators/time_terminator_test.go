package terminators

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimeTerminator(test *testing.T) {
	maximalDuration := 5 * time.Second
	terminator := NewTimeTerminator(clock, maximalDuration)

	gotClock := reflect.ValueOf(terminator.clock).Pointer()
	wantClock := reflect.ValueOf(clock).Pointer()
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

	for _, data := range []data{
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-4 * time.Second),
			},
			args: args{5},
			want: false,
		},
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-5 * time.Second),
			},
			args: args{5},
			want: true,
		},
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 5 * time.Second,
				startTime:       clock().Add(-6 * time.Second),
			},
			args: args{5},
			want: true,
		},
	} {
		terminator := TimeTerminator{
			clock:           data.fields.clock,
			maximalDuration: data.fields.maximalDuration,
			startTime:       data.fields.startTime,
		}
		got := terminator.IsSearchTerminated(data.args.deep)

		if got != data.want {
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

	for _, data := range []data{
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 100 * time.Second,
				startTime:       clock(),
			},
			args: args{5},
			want: 0,
		},
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-75 * time.Second),
			},
			args: args{5},
			want: 0.75,
		},
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-100 * time.Second),
			},
			args: args{5},
			want: 1,
		},
		{
			fields: fields{
				clock:           clock,
				maximalDuration: 100 * time.Second,
				startTime:       clock().Add(-110 * time.Second),
			},
			args: args{5},
			want: 1,
		},
	} {
		terminator := TimeTerminator{
			clock:           data.fields.clock,
			maximalDuration: data.fields.maximalDuration,
			startTime:       data.fields.startTime,
		}
		got := terminator.SearchProgress(data.args.deep)

		if got != data.want {
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
