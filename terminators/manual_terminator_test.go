package terminators

import (
	"testing"
)

func TestManualTerminatorIsSearchTerminated(test *testing.T) {
	type fields struct {
		terminationFlag uint64
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
			fields: fields{0},
			args:   args{5},
			want:   false,
		},
		{
			fields: fields{1},
			args:   args{5},
			want:   true,
		},
	} {
		terminator := ManualTerminator{
			terminationFlag: data.fields.terminationFlag,
		}
		got := terminator.IsSearchTerminated(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestManualTerminatorSearchProgress(test *testing.T) {
	type fields struct {
		terminationFlag uint64
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
			fields: fields{0},
			args:   args{5},
			want:   0,
		},
		{
			fields: fields{1},
			args:   args{5},
			want:   1,
		},
	} {
		terminator := ManualTerminator{
			terminationFlag: data.fields.terminationFlag,
		}
		got := terminator.SearchProgress(data.args.deep)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestManualTerminatorTerminate(test *testing.T) {
	var terminator ManualTerminator
	terminator.Terminate()

	flag := terminator.terminationFlag
	if flag != 1 {
		test.Fail()
	}
}
