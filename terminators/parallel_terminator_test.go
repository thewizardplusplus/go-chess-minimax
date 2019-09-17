package terminators

import (
	"testing"
)

func TestParallelTerminatorIsSearchTerminate(
	test *testing.T,
) {
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
		data{
			fields: fields{0},
			args:   args{5},
			want:   false,
		},
		data{
			fields: fields{1},
			args:   args{5},
			want:   true,
		},
	} {
		terminator := ParallelTerminator{
			terminationFlag: data.fields.
				terminationFlag,
		}
		got := terminator.IsSearchTerminate(
			data.args.deep,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
